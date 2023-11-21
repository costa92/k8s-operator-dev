# application-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/application-operator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/application-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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


## 测试
1. 运行 make envtest 
生成运行的二进制文件
![binary-file.png](docs%2Fimags%2Fbinary-file.png)
2. 修改 controllers/suite_test.go 中的 BeforeSuite 函数
![update-BeforSuite.png](docs%2Fimags%2Fupdate-BeforSuite.png)

```go
    // By default, tests run serially in the same process. To run multiple
	Expect(os.Setenv("TEST_ASSET_KUBE_APISERVER", "../bin/k8s/1.23.1-linux-amd64/kube-apiserver")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_ETCD", "../bin/k8s/1.23.1-linux-amd64/etcd")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_KUBECTL", "../bin/k8s/1.23.1-linux-amd64/kubectl")).To(Succeed())
```
对应的相关的环境变量

|   Variable name  |  	Type   |  	When to use   |     |
|-----|-----|-----|-----|
|   USE_EXISTING_CLUSTER  |  boolean   |   Instead of setting up a local control plane, point to the control plane of an existing cluster.  |     |
|   KUBEBUILDER_ASSETS  |  	path to directory   |   Point integration tests to a directory containing all binaries (api-server, etcd and kubectl).  |     |
|   TEST_ASSET_KUBE_APISERVER, TEST_ASSET_ETCD, TEST_ASSET_KUBECTL  |  	paths to, respectively, api-server, etcd and kubectl binaries   |   Similar to KUBEBUILDER_ASSETS, but more granular. Point integration tests to use binaries other than the default ones. <br/>These environment variables can also be used to ensure specific tests run with expected versions of these binaries.  |     |
|KUBEBUILDER_CONTROLPLANE_START_TIMEOUT and KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT | durations in format supported by time.ParseDuration | Set the timeout for starting and stopping the control plane. |  |
|KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT | boolean | If set to true, the output of the control plane will be attached to the test output. |  |

