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

// 检查pod是否存在
func (r *WorkSpaceReconciler) checkPodExist(key client.ObjectKey) (bool, error) {
	pod := &corev1.Pod{}
	// 先查询一下
	err := r.Client.Get(context.Background(), key, pod)
	if err != nil {
		// 判断是否存在
		if errors.IsNotFound(err) {
			return false, err
		}
		klog.Errorf("get pod error:%v", err)
		return false, err
	}
	return true, nil
}

// 删除 pod
func (r *WorkSpaceReconciler) deletePod(key client.ObjectKey) error {
	exist, err := r.checkPodExist(key)
	if err != nil {
		return err
	}
	// pod 不存在，直接返回
	if !exist {
		return nil
	}
	pod := &corev1.Pod{}
	pod.Name = key.Name           // 名字
	pod.Namespace = key.Namespace // 空间

	ctx, cancelFuc := context.WithTimeout(context.Background(), time.Second*35)
	defer cancelFuc()

	// 删除
	err = r.Client.Delete(ctx, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("delete pod error:%v", err)
		return err
	}
	return nil
}

func (r *WorkSpaceReconciler) constructPod(space *v1.WorkSpace) *corev1.Pod {
	volumeName := "volume-user-workspace"
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      space.Name,
			Namespace: space.Namespace,
			Labels: map[string]string{
				"app": "cloud-ide",
			},
		},

		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: space.Name,
							ReadOnly:  false,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            space.Name,
					Image:           space.Spec.Image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: space.Spec.Port,
						},
					},
					// 容器挂载储存卷
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      volumeName,
							ReadOnly:  false,
							MountPath: space.Spec.MountPath,
						},
					},
				},
			},
		},
	}

	if Mode == ModelRelease {

		pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},

			Limits: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(space.Spec.Cpu),
				corev1.ResourceMemory: resource.MustParse(space.Spec.Memory),
			},
		}
	}
	return pod
}
func (r *WorkSpaceReconciler) createPod(space *v1.WorkSpace, key client.ObjectKey) error {
	// 1.检查Pod是否存在
	exist, err := r.checkPodExist(key)
	if err != nil {
		return err
	}

	// Pod已存在,直接返回
	if exist {
		return nil
	}

	// 2.创建Pod
	pod := r.constructPod(space)

	// 设置控制器，如果设置了控制器,那么被控制的资源的变化也会被发送到队列中
	//if err = controllerutil.SetControllerReference(space, pod, r.Scheme); err != nil {
	//	return err
	//}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	err = r.Client.Create(ctx, pod)
	if err != nil {
		// 如果Pod已经存在,直接返回
		if errors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}
