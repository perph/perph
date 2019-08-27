/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	agent "github.com/perph/perph/api/v1"
	"github.com/perph/perph/configure"
)

// CheckReconciler reconciles a Check object
type CheckReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=agent.perph.io,resources=checks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=agent.perph.io,resources=checks/status,verbs=get;update;patch

func (r *CheckReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("check", req.NamespacedName)

	//Get Check Config
	var check agent.Check
	if err := r.Get(ctx, req.NamespacedName, &check); err != nil {
		log.Error(err, "unable to fetch the check config")
		return ctrl.Result{}, ignoreNotFound(err)
	}
	//Get Associated Endpoints
	var endpoints corev1.Endpoints
	if err := r.Get(ctx, types.NamespacedName{Name: check.ObjectMeta.Name, Namespace: req.Namespace}, &endpoints); err != nil {
		log.Error(err, "unable to fetch the endpoints for the synthetic run agents")
		//Reschedule this after 5 seconds, the endpoints should hopefully be back up
		return ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}

	r.Log.Info("endpoints", "ep", endpoints)
	//make an RPC call to configure the agent
	addresses, notReadyAddr := getAddress(endpoints)
	//TODO configure this port to something more dynanmic later
	for _, a := range addresses {
		a = a + ":12000"
		r.Log.Info(a)
		_, err := configure.ConfigureAgent(a, check.Spec.JobID)
		if err != nil {
			r.Log.Error(err, "unable to notify agent of updated configuration - retrying in 5 seconds")
			return ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}
	}
	if len(notReadyAddr) != 0 {
		r.Log.Info("some endpoints are not ready - retrying in 5 seconds")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	r.Log.Info(check.Spec.JobID, "JOB ID CONFIGURED")
	//If error reschedule the attempt in 5 seconds

	return ctrl.Result{}, nil

}

func (r *CheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agent.Check{}).
		Watches(&source.Kind{Type: &corev1.Endpoints{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []ctrl.Request {
				//Check if it was caused by a change in endpoints
				var ep corev1.Endpoints
				err := r.Get(context.Background(), types.NamespacedName{Name: obj.Meta.GetName(), Namespace: obj.Meta.GetNamespace()}, &ep)
				if err != nil {
					r.Log.Info("the event was not a change in endpoints", "obj", obj)
					return nil
				}

				//Filter Endpoints to those associated with a running synthetic run job
				jobID, ok := ep.ObjectMeta.Labels["job"]
				if !ok {
					// r.Log.Info("the endpoints are not labelled with a jobID", "obj", ep)
					return nil
				}

				lbs := map[string]string{"perph": "agent", "job": jobID}
				if !ContainsLabels(ep.ObjectMeta.Labels, lbs) {
					return nil
				}

				r.Log.Info("the endpoints are", "obj", ep)
				// Add job name and namespace into the request
				res := make([]ctrl.Request, 1)
				res[0].Name = ep.Name
				res[0].Namespace = ep.Namespace
				return res
			}),
		}).
		Complete(r)
}

//ContainsLabels -  Checks that all labels are contained within the supplied object label map
func ContainsLabels(objLabels map[string]string, labels map[string]string) bool {
	i := 0
	for key, label := range labels {
		if objLabels[key] == label {
			i++
		}
	}
	if i != len(labels) {
		return false
	}
	return true
}

func getAddress(ep corev1.Endpoints) (addr []string, notReadyAddr []string) {
	for _, sub := range ep.Subsets {
		for _, a := range sub.Addresses {
			addr = append(addr, a.IP)
		}
		for _, a := range sub.NotReadyAddresses {
			notReadyAddr = append(notReadyAddr, a.IP)
		}
	}
	return
}
