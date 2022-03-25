/*
Copyright 2022.

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
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	samplevolumnev1 "demo-volume/api/v1"
)

// SampleVolumeReconciler reconciles a SampleVolume object
type SampleVolumeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=samplevolumne.operator.yogeshsharma.me,resources=samplevolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=samplevolumne.operator.yogeshsharma.me,resources=samplevolumes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=samplevolumne.operator.yogeshsharma.me,resources=samplevolumes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SampleVolume object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SampleVolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Enter reconcile", "req", req)
	volume := &samplevolumnev1.SampleVolume{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, volume)

	logger.Info("Enter reconcile", "spec", volume.Spec, "status", volume.Status)
	if volume.Spec.Name != volume.Status.Name {
		volume.Status.Name = volume.Spec.Name
		r.Status().Update(ctx, volume)
	}

	r.reconcilePVC(ctx, volume, logger)
	return ctrl.Result{}, nil
}

func (r *SampleVolumeReconciler) reconcilePVC(ctx context.Context, volume *samplevolumnev1.SampleVolume, logger logr.Logger) error {
	pvc := &v1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: volume.Name, Namespace: volume.Namespace}, pvc)
	if err == nil {
		logger.Info("PVC found in namespace", "namespace", volume.Namespace)
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	logger.Info("PVC Not found in namespace", "namespace", volume.Namespace)
	//storageClass := "yogeshsharma-volume"
	storageRequest, _ := resource.ParseQuantity(fmt.Sprintf("%dGi", volume.Spec.Size))

	pvc = &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volume.Name,
			Namespace: volume.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			//StorageClassName: &storageClass,
			AccessModes: []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{"storage": storageRequest},
			},
		},
	}

	logger.Info("Status of the volume is", "status", volume.Status)

	logger.Info("Setting up pvc controller reference")
	if err = controllerutil.SetControllerReference(volume, pvc, r.Scheme); err != nil {
		logger.Error(err, "Failed to set pvc controller reference")
		return err
	}

	return r.Create(ctx, pvc)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SampleVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplevolumnev1.SampleVolume{}).
		Complete(r)
}
