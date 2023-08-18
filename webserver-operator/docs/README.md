# webserver-operator docs

## install kubebuilder

```bash
$git clone git@github.com:kubernetes-sigs/kubebuilder.git
$cd kubebuilder
$make && cd bin 
$./kubebuilder version
$sudo cp kubebuilder /usr/local/bin 
$kubebuilder version
```

## create project

```bash
$kubebuilder init --repo=github.com/costa92/webserver-operator --project-name webserver-operator
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.14.1
Update dependencies:
$ go mod tidy
Next: define a resource with:
$ kubebuilder create api
```
 > 注：–repo指定go.mod中的module root path，你可以定义你自己的module root path。

## create api,init CRD

```bash
kubebuilder create api --version v1 --kind WebServer 
Create Resource [y/n]
y
Create Controller [y/n]
y
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1/webserver_types.go
controllers/webserver_controller.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
mkdir -p /home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin
test -s /home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin/controller-gen && /home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin/controller-g
en --version | grep -q v0.11.1 || \
GOBIN=/home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.1
go: downloading sigs.k8s.io/controller-tools v0.11.1
go: downloading golang.org/x/tools v0.4.0
go: downloading github.com/gobuffalo/flect v0.3.0
go: downloading golang.org/x/mod v0.7.0
/home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

### 执行 make manifests 生成CRD的yaml文件
```bash
$make manifests
/home/hellotalk/code/go/src/github.com/costa92/operation/webserver-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

## 查看项目结构
```bash
$tree -F .
```

## webserver-operator的基本结构

![webserver-operator的基本结构](https://tonybai.com/wp-content/uploads/developing-kubernetes-operators-in-go-part1-7.png)


### 说明

1. CRD
代码中生成的 CRD yaml 文件位于 config/crd/bases/my.domain_webservers.yaml，CRD 与 api/v1/webserver_types.go 密切相关
,在 api/v1/webserver_types.go 中为CRD定义spec 相关字段，之后 make manifests命令可以解析webserver_types.go中的变化并更新CRD的yaml文件
2. Controller
代码中生成的 Controller 自身就是作为一个 Deployment 运行在 Kubernetes 集群中，是监视 CRD 运行状态，根据 CRD 的状态变化执行相应的操作，
比如创建一个 Deployment，或者删除一个 Deployment。 并在 Reconcile 方法中实现了对 CRD 的状态变化的监视和处理逻辑。 
3. 其他
  权限控制：controller 通过 serviceAccount 访问 k8s API Service,通过 config/rbac/role.yaml 和 config/rbac/role_binding.yaml 为 serviceAccount 分配权限
