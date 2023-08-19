package controllers

import (
	"context"
	v1 "github.com/costa92/cloud-ide-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func (r *WorkSpaceReconciler) checkPVCExist(key client.ObjectKey) (bool, error) {
	pvc := &corev1.PersistentVolumeClaim{}

	if err := r.Client.Get(context.Background(), key, pvc); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		klog.Errorf("get pvc error:%v", err)
		return false, nil
	}
	return true, nil
}

func (r *WorkSpaceReconciler) deletePVC(key client.ObjectKey) error {
	exist, err := r.checkPVCExist(key)
	if err != nil {
		return err
	}

	// pvc 不存在,无需删除
	if !exist {
		return nil
	}
	pvc := &corev1.PersistentVolumeClaim{}
	pvc.Name = key.Name           // 名字
	pvc.Namespace = key.Namespace // 空间

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	// 删除
	if err := r.Client.Delete(ctx, pvc); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("delete pvc error:%v", err)
		return err
	}
	return nil
}

func (r *WorkSpaceReconciler) createPVC(space *v1.WorkSpace, key client.ObjectKey) error {
	//  1 先检查
	exist, err := r.checkPodExist(key)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	// 2 创建 pvc
	pvc, err := r.constructPVC(space)
	if err != nil {
		klog.Errorf("construct pvc error:%v", err)
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	err = r.Client.Create(ctx, pvc)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}

	return nil
}

func (r *WorkSpaceReconciler) constructPVC(space *v1.WorkSpace) (*corev1.PersistentVolumeClaim, error) {
	quantity, err := resource.ParseQuantity(space.Spec.Storage)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{corev1.ResourceStorage: quantity},
				Requests: corev1.ResourceList{corev1.ResourceStorage: quantity},
			},
		},
	}
	return pvc, nil
}
