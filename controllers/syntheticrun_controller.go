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
	"fmt"
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	agent "github.com/perph/perph/api/v1"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SyntheticRunReconciler reconciles a SyntheticRun object
type SyntheticRunReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=agent.perph.io,resources=syntheticruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=agent.perph.io,resources=syntheticruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;create;update;watch;list

func (r *SyntheticRunReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("syntheticrun", req.NamespacedName)

	log.Info("starting reconciliation")
	//Get Synthetic Runs
	var synthRun agent.SyntheticRun
	if err := r.Get(ctx, req.NamespacedName, &synthRun); err != nil {
		log.Error(err, "unable to fetch the synthetic run agent")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	//Check if a deployment needs creating or updating
	if _, err := r.deployAgent(ctx, &synthRun); err != nil {
		log.Error(err, "unable to create deployment for synthetic run")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	//The synthRun.Spec.JobID will be replaced with a compare of config CRD Spec.JobID version
	//The config jobID will be added in via mutating admission webhook
	if synthRun.Spec.JobID != synthRun.Status.JobID {
		//Get a list of pods
		var agentInstances corev1.PodList
		if err := r.List(ctx, &agentInstances, client.InNamespace(req.Namespace), client.MatchingLabels(map[string]string{"perph": "agent", "job": synthRun.Name})); err != nil {
			log.Error(err, "unable to fetch the synthetic run agent instances")
			return ctrl.Result{}, err
		}

		//make an RPC call to configure the agent
		//TODO add an RPC here
		err := dummyGRPC(synthRun.Spec.JobID)
		//If error reschedule the attempt in 5 seconds
		if err != nil {
			r.Log.Error(err, "unable to notify agent of updated configuration - retrying in 5 seconds")
			return ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}

		synthRun.Status.JobID = synthRun.Spec.JobID

		if err := r.Status().Update(ctx, &synthRun); err != nil {
			log.Error(err, "unable to update syntheticagent status")
		}

	}
	log.Info("success")

	return ctrl.Result{}, nil
}

func (r *SyntheticRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agent.SyntheticRun{}).
		Owns(&apps.Deployment{}).
		Complete(r)
}

func (r *SyntheticRunReconciler) deployAgent(ctx context.Context, synthRun *agent.SyntheticRun) (*apps.Deployment, error) {
	//Set the deployment size based on spec or default
	requiredInstances := synthRun.Spec.InstanceCount
	if requiredInstances == 0 {
		requiredInstances = 3
	}

	//Create deployment
	depl := &apps.Deployment{ObjectMeta: metav1.ObjectMeta{Name: synthRun.Name + "-deploy", Namespace: synthRun.Namespace}}
	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, depl, func() error {
		lbs := map[string]string{"perph": "agent", "job": synthRun.Name}
		depl.Spec.Replicas = &requiredInstances

		depl.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: lbs,
		}

		depl.Spec.Template.ObjectMeta.Labels = lbs
		depl.Spec.Template.Spec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "agent",
					Image: "perph/agent",
					Ports: []corev1.ContainerPort{{ContainerPort: 6379}},
				},
			},
		}

		if err := ctrl.SetControllerReference(synthRun, depl, r.Scheme); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.Log.Error(err, "unable to ensure deployment is up to date", "agent", synthRun.Name)
		if err := r.Status().Update(ctx, synthRun); err != nil {
			r.Log.Error(err, "unable to update agent status")
		}
		return nil, err
	}
	return depl, nil
}

func dummyGRPC(id string) error {
	n := rand.Intn(10)
	if n > 7 {
		return fmt.Errorf("failed RPC randomly")
	}
	return nil
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
