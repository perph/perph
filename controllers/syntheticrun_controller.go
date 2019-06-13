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

	syntheticv1 "github.com/perph/perph/api/v1"
)

// SyntheticRunReconciler reconciles a SyntheticRun object
type SyntheticRunReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=synthetic.perph.io,resources=syntheticruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=synthetic.perph.io,resources=syntheticruns/status,verbs=get;update;patch

func (r *SyntheticRunReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("syntheticrun", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *SyntheticRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syntheticv1.SyntheticRun{}).
		Owns(&batch.Job{}).
		// Watches(&source.Kind{Type: &syntheticv1.Validation{}}, &handler.EnqueueRequestsFromMapFunc{
		// 	// ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []ctrl.Request {
		// 	// 	var runs syntheticv1.SyntheticRun
		// 	// 	if err := r.List(context.Background(), &runs, client.InNamespace(obj.Meta.GetNamespace()), client.MatchingField(".spec.validationName", obj.Meta.GetName())); err != nil {
		// 	// 		r.Log.Info("unable to get synthentic run for validation", "validation", obj)
		// 	// 		return nil
		// 	// 	}

		// 	// 	res := make([]ctrl.Request, len(runs.Validations))
		// 	// 	for i, run := range runs.Validations {
		// 	// 		res[i].Name = run.Name
		// 	// 		res[i].Namespace = run.Namespace
		// 	// 	}
		// 	// 	return res
		// 	// }),
		// }).
		// Watches(&source.Kind{Type: &syntheticv1.Check{}}, &handler.EnqueueRequestForObject{
		// 	//TODO add
		// 	ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []ctrl.Request {
		// 		var runs syntheticv1.SyntheticRun
		// 		if err := r.List(context.Background(), &runs, client.InNamespace(obj.Meta.GetNamespace()), client.MatchingField(".spec.checkName", obj.Meta.GetName())); err != nil {
		// 			r.Log.Info("unable to get synthentic run for check", "check", obj)
		// 			return nil
		// 		}

		// 		res := make([]ctrl.Request, len(runs.Checks))
		// 		for i, run := range runs.Checks {
		// 			res[i].Name = run.Name
		// 			res[i].Namespace = run.Namespace
		// 		}
		// 		return res
		// 	}),
		// }).
		Complete(r)
}
