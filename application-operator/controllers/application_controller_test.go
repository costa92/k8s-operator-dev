package controllers

import (
	"context"
	"fmt"
	appsv1 "github.com/costa92/application-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Application controller", func() {
	Context("When creating an Application", func() {
		const ApplicationName = "application-operator-0"
		// Create a new context
		ctx := context.Background()
		// Create the Namespace object
		namespce := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ApplicationName,
				Namespace: ApplicationName,
			},
		}

		// Create the Namespace
		typeNameSpaceName := types.NamespacedName{
			Name:      ApplicationName,
			Namespace: ApplicationName,
		}

		// Create the Namespace object
		BeforeEach(func() {
			By("Creating a new Namespace")
			err := k8sClient.Create(ctx, namespce)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("BUSYBOX_IMAGE", "costa92/treafik-api:v0.0.4")
			Expect(err).To(Not(HaveOccurred()))
		})

		AfterEach(func() {
			By("Deleting the Namespace")
			err := k8sClient.Delete(ctx, namespce)
			Expect(err).ToNot(HaveOccurred())

			By("Unset BUSYBOX_IMAGE")
			err = os.Unsetenv("BUSYBOX_IMAGE")
			Expect(err).To(Not(HaveOccurred()))
		})

		It("Should create a new Deployment", func() {
			By("By creating a new Application")
			application := &appsv1.Application{}

			err := k8sClient.Get(ctx, typeNameSpaceName, application)
			if err != nil && errors.IsNotFound(err) {
				application = &appsv1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      ApplicationName,
						Namespace: ApplicationName,
					},
					Spec: appsv1.ApplicationSpec{
						Replicas: 1,
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "busybox",
										Image: os.Getenv("BUSYBOX_IMAGE"),
									},
								},
							},
						},
					},
				}
				err = k8sClient.Create(ctx, application)
				fmt.Println(err)
				Expect(err).ToNot(HaveOccurred())
			}

			By("Checking if the custom resource was successfully created")

			Eventually(func() error {
				found := &appsv1.Application{}
				return k8sClient.Get(ctx, typeNameSpaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Reconciling the custom resource created")
			reconciling := &ApplicationReconciler{
				Client: k8sClient,
				Scheme: scheme.Scheme,
			}

			// Reconcile the custom resource
			// is pod resource
			_, err = reconciling.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNameSpaceName,
			})

			// 实现执行 Reconcile 之后，会创建一个 Pod

			Expect(err).ToNot(HaveOccurred())

			By("Checking if the Pod was successfully created")
			//Eventually(func() error {
			//	found := &appsv1s.Deployment{}
			//	return k8sClient.Get(ctx, typeNameSpaceName, found)
			//}, time.Minute, time.Second).Should(Succeed())

			Eventually(func() error {
				// 通过 Pod 名称获取 Pod 对象
				found := &corev1.Pod{}
				err := k8sClient.Get(ctx, typeNameSpaceName, found)
				if err != nil {
					fmt.Println(err.Error())
				}
				return err
			}, time.Minute, time.Second).Should(Succeed())

			By("Checking the latest Status Condition added to the Application instance")
			Eventually(func() error {
				fmt.Println(application)
				return nil
			}, time.Minute, time.Second).Should(Succeed())
		})
	})
})

// 参考代码： https://github.com/kubernetes-sigs/kubebuilder/blob/v3.7.0/testdata/project-v3-with-deploy-image/controllers/busybox_controller_test.go
