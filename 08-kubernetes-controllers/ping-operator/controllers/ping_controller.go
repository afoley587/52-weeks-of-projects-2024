/*
Copyright 2024.

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
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitorsv1beta1 "github.com/afoley587/52-weeks-of-projects-2023/08-kubernetes-controllers/ping-operator/api/v1beta1"
)

// PingReconciler reconciles a Ping object
type PingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitors.engineeringwithalex.io,resources=pings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitors.engineeringwithalex.io,resources=pings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitors.engineeringwithalex.io,resources=pings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ping object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var ping monitorsv1beta1.Ping

	if err := r.Get(ctx, req.NamespacedName, &ping); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	job, err := r.BuildJob(ping)

	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.Create(ctx, &job); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

func (r *PingReconciler) BuildJob(ping monitorsv1beta1.Ping) (batchv1.Job, error) {
	attempts := "-c" + strconv.Itoa(ping.Spec.Attempts)
	host := ping.Spec.Hostname
	j := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: batchv1.SchemeGroupVersion.String(),
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ping.Name + "-job",
			Namespace: ping.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "ping",
							Image:   "bash",
							Command: []string{"/bin/ping"},
							Args:    []string{attempts, host},
						},
					},
				},
			},
		},
	}
	return j, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorsv1beta1.Ping{}).
		Complete(r)
}
