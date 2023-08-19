/*
Copyright 2023 Costalong.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/costa92/cloud-ide-operator/api/v1"
)

var Mode string

const (
	ModelRelease = "release"
	ModDev       = "dev"
)

// WorkSpaceReconciler reconciles a WorkSpace object
type WorkSpaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.costalong.com,resources=workspaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.costalong.com,resources=workspaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.costalong.com,resources=workspaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WorkSpace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
// Reconcile的意思是协调
func (r *WorkSpaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// 先查询 WorkSpace
	wp := appsv1.WorkSpace{}
	err := r.Client.Get(context.Background(), req.NamespacedName, &wp)

	// case 1  没有查到 Workspace，说明 WorkSpace 被删除了，删除对应的Pod 和 PVC 即可
	if err != nil {
		if errors.IsNotFound(err) {
			if e1 := r.deletePod(req.NamespacedName); e1 != nil {
				klog.Errorf("[Delete Workspace] delete pod error:%v", e1)
				return ctrl.Result{Requeue: true}, e1
			}

			if e2 := r.deletePVC(req.NamespacedName); e2 != nil {
				klog.Errorf("[Delete Workspace] delete pvc error:%v", e2)
				return ctrl.Result{Requeue: true}, e2
			}
			return ctrl.Result{}, nil
		}
		klog.Errorf("get workspace error:%v", err)
		return ctrl.Result{Requeue: true}, err
	}

	// 找到了 WorkSpace,根据 WorkSpace 的operation 字段 判断进行操作
	switch wp.Spec.Operation {
	// case 2: 启动 workspace 检查 pvc 是否存在
	case appsv1.WorkSpaceStart:
		err := r.createPVC(&wp, req.NamespacedName)
		if err != nil {
			klog.Errorf("[start Workspace] create pvc error:%v", err)
			return ctrl.Result{Requeue: true}, err
		}

		// 创建Pod
		err = r.createPod(&wp, req.NamespacedName)
		if err != nil {
			klog.Errorf("[Start Workspace] create pod error:%v", err)
			return ctrl.Result{Requeue: true}, err
		}
		r.updateStatus(&wp, appsv1.WorkspacePhaseRunning)
	case appsv1.WorkSpaceStop:
		// 删除 pod
		err = r.deletePod(req.NamespacedName)
		if err != nil {
			klog.Errorf("[Stop workspace] delete pod error:%v", err)
			return ctrl.Result{Requeue: true}, err
		}
		r.updateStatus(&wp, appsv1.WorkspacePhaseStopped)
	}
	return ctrl.Result{}, nil
}

func (r WorkSpaceReconciler) updateStatus(wp *appsv1.WorkSpace, phase appsv1.WorkSpacePhase) {
	wp.Status.Phase = phase
	err := r.Client.Status().Update(context.Background(), wp)
	if err != nil {
		klog.Errorf("update status error:%v", err)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkSpaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 8}).
		For(&appsv1.WorkSpace{}).
		Owns(&corev1.Pod{}, builder.WithPredicates(predicatePod)).
		Owns(&corev1.PersistentVolumeClaim{}, builder.WithPredicates(predicatePVC)).
		Complete(r)
}
