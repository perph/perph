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
	agent "github.com/perph/perph/api/v1"
	"github.com/perph/perph/configure"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

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

	//Check if a svc is deployed or create one
	depl, err := r.deployService(ctx, &synthRun)
	if err != nil {
		log.Error(err, "unable to create svc for synthetic run")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	//The synthRun.Spec.JobID will be replaced with a compare of config CRD Spec.JobID version
	//The config jobID will be added in via mutating admission webhook
	if synthRun.Spec.JobID != synthRun.Status.JobID {
		//Get a list of eps
		var endpoints corev1.Endpoints
		if err := r.Get(ctx, types.NamespacedName{Name: depl.ObjectMeta.Name, Namespace: req.Namespace}, &endpoints); err != nil {
			log.Error(err, "unable to fetch the endpoints for the synthetic run agentss")
			return ctrl.Result{}, err
		}

		//make an RPC call to configure the agent
		addresses := getAddress(endpoints)
		address := addresses[0] + ":12000"
		r.Log.Info(address)
		_, err := configure.ConfigureAgent(address, synthRun.Spec.JobID)
		// r.Log.Info(jobId, "JOB ID CONFIGURED")
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
		Owns(&corev1.Service{}).
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
					Name:            "agent",
					Image:           "perph.io/agent:latest",
					ImagePullPolicy: "IfNotPresent",
					// Ports:           []corev1.ContainerPort{corev1.ContainerPort{ContainerPort: 12000}},
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

func (r *SyntheticRunReconciler) deployService(ctx context.Context, synthRun *agent.SyntheticRun) (*corev1.Service, error) {
	//Create deployment
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: synthRun.Name + "-svc", Namespace: synthRun.Namespace}}
	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, svc, func() error {
		lbs := map[string]string{"perph": "agent", "job": synthRun.Name}
		svc.Spec.Selector = lbs
		svc.Spec.Type = "ClusterIP"
		svc.Spec.Ports = []corev1.ServicePort{{Name: "agents", Port: 12000, TargetPort: intstr.FromInt(12000)}}

		if err := ctrl.SetControllerReference(synthRun, svc, r.Scheme); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.Log.Error(err, "unable to ensure svc is up to date", "agent", synthRun.Name)
		if err := r.Status().Update(ctx, synthRun); err != nil {
			r.Log.Error(err, "unable to update agent status")
		}
		return nil, err
	}
	return svc, nil
}

func getAddress(ep corev1.Endpoints) []string {
	var addresses []string
	for _, sub := range ep.Subsets {
		for _, addr := range sub.Addresses {
			addresses = append(addresses, addr.IP)
		}
	}
	return addresses
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
