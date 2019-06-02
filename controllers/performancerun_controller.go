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

	"github.com/go-logr/logr"
	batch "k8s.io/api/batch/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"

	syntheticv1 "github.com/perph/perph/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

// PerformanceRunReconciler reconciles a PerformanceRun object
type PerformanceRunReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=synthetic.perph.io,resources=performanceruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=synthetic.perph.io,resources=performanceruns/status,verbs=get;update;patch

func (r *PerformanceRunReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("performancerun", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *PerformanceRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syntheticv1.PerformanceRun{}).
		Owns(&batch.Job{}).
		Watches(&source.Kind{Type: &syntheticv1.LoadTest{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []ctrl.Request {
				var runs syntheticv1.PerformanceRun
				if err := r.List(context.Background(), &runs, client.InNamespace(obj.Meta.GetNamespace()), client.MatchingField(".spec.loadtestName", obj.Meta.GetName())); err != nil {
					r.Log.Info("unable to get performance run for loadtest", "loadtest", obj)
					return nil
				}

				res := make([]ctrl.Request, len(runs.Validations))
				for i, run := range runs.Validations {
					res[i].Name = run.Name
					res[i].Namespace = run.Namespace
				}
				return res
			}),
		}).
		Complete(r)
}
