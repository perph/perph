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
		Complete(r)
}
