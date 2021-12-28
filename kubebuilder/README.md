# Kubebuilder

kubebuilder是一个在go中快速构建和发布Kubernetes API的框架（使用CRDs），[简化了传统创建CRD并编写相应controller的流程](https://www.cnblogs.com/yrxing/p/14472018.html)

## 快速入门

安装

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.2.0/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.2.0_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

创建一个project

```bash
kubebuilder init --domain my.domain
```

创建一个API

```bash
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

执行完毕后，会在config文件夹生成相应的CRD的yaml文件

api文件夹是相应的类型定义，需要自定义信息的话需要修改此文件夹中的内容。如果修改了相应类型，config/samples中的CRD文件也要一起修改。

controllers是相应controller的定义模版，之后编写相应的业务逻辑是在这里。

生成的API的版本信息为：schema.GroupVersion{Group: "webapp.my.domain", Version: "v1"}



集群中安装/卸载CRD

```bash
make install/uninstall
```

运行controller

```go
make run
```



集群中运行Operator

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make deploy IMG=<some-registry>/<project-name>:tag
```

config/manager文件夹定义了这个API相关的部署文件，make deploy所做的工作就是替换掉deployment中的image，改为我们打包好的operator的image。



由上可见，kubebuilder可以让我们专注于controller的编写，省去了其他不必要的工作，比如API的定义，CRD的定义等。

## 从零编写一个CrontJob

现在让我们假设 Kubernetes 中实现的 `CronJob Controller` 并不能满足我们的需求，我们希望使用 Kubebuilder 重写它。

`CronJob controller` 会控制 kubernetes 集群上的 job 每隔一段时间运行一次，它是基于 `Job controller` 实现的，`Job controller` 的 job 只会执行任务一次。

首先创建项目

```bash
kubebuilder init --domain tutorial.kubebuilder.io
```



### 项目结构

用于构建和部署controller的Makefile文件。

用于构建镜像的Dockerfile文件。

用于创建新组件的元数据PROJECT文件。

一个config目录：包含运行operator所需要的所有配置文件，后续编写operator的时候还会包含CRD、RBAC、Webhook等相关配置文件。

[`config/default`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/default) 包含 [Kustomize base](https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/book/src/cronjob-tutorial/testdata/project/config/default/kustomization.yaml) 文件，用于以标准配置启动 controller。

[`config/`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config)目录下每个目录都包含不同的配置：

- [`config/manager`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/manager): 包含在 k8s 集群中以 pod 形式运行 controller 的 YAML 配置文件
- [`config/rbac`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/rbac): 包含运行 controller 所需最小权限的配置文件

程序入口main.go

现在main里面还没有很多东西，首先实例化一个scheme，每组controller都需要一个scheme，它提供Kinds与Go types之间的银蛇关系。

然后为metric绑定了一些基本的flages

实例化一个controller-runtime库中的manager，它用于跟踪所有的controller，manager一旦运行就会一直运行下去直到收到正常的关闭信号。这样当operator运行在kubernetes上时，可以通过优雅的方式终止这个Pod。

虽然目前我们还没有什么可以运行，但是请记住 `+kubebuilder:scaffold:builder` 注释的位置 -- 很快那里就会变得有趣。



### 创建API

Kind和Resouce的区别：

Resouce是Kind在API中的表示，通常情况下是一一对应的，但有时候一个Kind可能有多个Resources。比如 Scale Kind 可能对应很多 Resources：deployments/scale 或者 replicasets/scale, 但是在 CRD 中，每个 `Kind` 只会对应一种 `Resource`。

当我们使用 kubectl 操作 API 时，操作的就是 `Resource`，比如 `kubectl get pods`, 这里的 `pods` 就是指 `Resource`。

而我们在编写 YAML 文件时，会编写类似 `Kind: Pod` 这样的内容，这里 `Pod` 就是 `Kind`

稍后我们将看到每个 `GVK` 都对应一个 root Go type (比如：Deployment 就关联着 K8s 源码里面 k8s.io/api/apps/v1 package 中的 Deployment struct)。

Scheme就提供了GVK和Go type之间的映射，这种关系通过Scheme注册到API Server中后，我们就可以从中获取数据并反序列化到相应结构体中了。

接下来开始使用kubebuilder创建一个API：

```
kubebuilder create api --group batch --version v1 --kind CronJob
```

api/v1目录会被创建，这个目录有三个文件：crontjob_types.go/groupversion_info.go/zz_geneated.deepcopy.go。

每次我们运行这个命令但是指定不同的 `Kind` 时, 他都会为我们创建一个 `xxx_types.go` 文件。



解析crontjob_types.go文件：

这个文件定一个CrontJob的期望状态Spec和实际状态Status，以及对应的结构体类型CrontJob和CrontJobList。

`CronJob` 是我们的 root type，用来描述 `CronJob Kind`。和所有 Kubernetes 对象一样， 它包含 `TypeMeta` (用来定义 API version 和 Kind) 和 `ObjectMeta` (用来定义 name、namespace 和 labels等一些字段)

`CronJobList` 包含了一个 `CronJob` 的切片，它是用来批量操作 `Kind` 的，比如 LIST 操作。通常不会修改这两个，所有的修改都是在Spec和Status上进行的。

`+kubebuilder:object:root` 注释称为标记(marker)。 稍后我们还会看到更多它们，它们提供了一些元数据， 来告诉 [controller-tools](https://github.com/kubernetes-sigs/controller-tools)(我们的代码和 YAML 生成器) 一些额外的信息。 这个注释告诉 `object` 这是一种 `root type Kind`。 然后，<font color=red>`object` 生成器会为我们生成 [runtime.Object](https://godoc.org/k8s.io/apimachinery/pkg/runtime#Object) 接口的实现， 这是所有 Kinds 必须实现的接口。</font>

```go
func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}
```

最后将Kinds注册到API group中。



groupversion_info.go文件：

首先，我们有一些 `package-level` 的标记，`+kubebuilder:object:generate=true` 表示该程序包中有 Kubernetes 对象， `+groupName=batch.tutorial.kubebuilder.io` 表示该程序包的 Api group 是 `batch.tutorial.kubebuilder.io`。 `object` 生成器使用前者，而 CRD 生成器使用后者, 来为由此包创建的 CRD 生成正确的元数据。

然后定义了一些全局变量帮助我们建立Scheme。



zz_generated.deepcopy.go

`zz_generated.deepcopy.go` 包含[之前所说的](https://xuejipeng.github.io/kubebuilder-doc-cn/cronjob-tutorial/new-api.html#moment)由 `+kubebuilder:object:root`自动生成的 `runtime.Object` 接口的实现。



### 设计一个API

在kubernetes中，有一些设计API的规范:

- 所有可序列化字段必须是camelCase格式，需要用json tag字段指定。

- 字段为空时，用omitempty标记。

- 整数只能使用int32和int64，小数只能使用resouce.Quantity

  Quantity 是十进制数字的一种特殊表示法，具有明确固定的表示形式， 使它们在计算机之间更易于移植。在Kubernetes中指定资源请求和Pod的限制时，您可能已经注意到它们。例如，该值2m表示0.002十进制表示法。 2Ki 表示2048十进制，而2K表示2000十进制。

- 时间使用metav1.Time,它除了格式在 Kubernetes 中比较通用外，功能与 `time.Time` 完全相同，可以获得更稳定的序列化。

具体CrontJob的样子，参考文件。

一些mark的功能：

- +kubebuilder:validation:Enum=Allow;Forbid;Replace：这个类型只接受这三个值，validation用来对类型做一些合法性检查。
- +kubebuilder:subresouce:status：定义在Kind结构前面，如果希望操作该Kind像kubernetes内置资源一样，增加这个mark。



### 实现controller

每个controller专注于一个root kind，但也可以与其他Kinds进行交互。

controller的职责就是确保实际状态与期望状态相匹配，匹配过程称为reconciling（调和，使...一直）。

在 controller-runtime 库中，实现 Kind reconciling 的逻辑我们称为 [*Reconciler*](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile)。 `reconciler` 获取对象的名称并返回是否需要重试(例如: 发生错误或是一些周期性的 controllers，像 HorizontalPodAutoscale)。

大多数 controllers 最终都会运行在 k8s 集群上，因此它们需要 RBAC 权限, 我们使用 controller-tools [RBAC markers](https://xuejipeng.github.io/reference/markers/rbac.html) 指定了这些权限。 这是运行所需的最小权限。 随着我们添加更多功能，我们将会重新定义这些权限。

kubebuilder 为我们搭建了一个基本的 `reconciler` 结构体。 几乎每个 `reconciler` 都需要记录日志，并且需要能够获取对象，因此这个结构体是开箱即用的。

```go
func (r *CronJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
   _ = context.Background()
   _ = r.Log.WithValues("cronjob", req.NamespacedName)

   // your logic here

   return ctrl.Result{}, nil
}

func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
   return ctrl.NewControllerManagedBy(mgr).
      For(&batchv1.CronJob{}).
      Complete(r)
}
```

另外还绑定了两个方法Reconcile和SetupWithManager，`Reconcile` 方法对个某个单一的 object 执行 `reconciling` 动作， 我们的 [Request](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Request)只是一个 name， 但是 client 可以通过 name 信息从 cache 中获取到对应的 object。

最后，我们将此 reconciler 加到 manager，以便在启动 manager 时启动 reconciler。



接下来就是编写我们自己controller的业务逻辑了，我们的 CronJob controller 的基本逻辑是：

1. 按namespace加载 CronJob
2. 列出所有 active jobs，并更新状态
3. 根据历史记录清理 old jobs
4. 检查 Job 是否已被 suspended（如果被 suspended，请不要执行任何操作）
5. 获取到下一次要 schedule 的 Job
6. 运行新的 Job, 确定新 Job 没有超过 deadline 时间，且不会被我们 concurrency 规则 block
7. 如果 Job 正在运行或者它应该下次运行，请重新排队



如果你想为你的 CRD 实现 [admission webhooks](https://xuejipeng.github.io/kubebuilder-doc-cn/reference/admission-webhook.html)，你只需要实现 `Defaulter` 和 (或) `Validator` 接口即可。

其余的东西 Kubebuilder 会为你实现，比如：

1. 创建一个 webhook server
2. 确保这个 server 被添加到 manager 中
3. 为你的 webhooks 创建一个 handlers
4. 将每个 handler 以 path 形式注册到你的 server 中

首先，创建webhooks

```
 kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```

这将创建 Webhook 功能相关的方法，并在 `main.go` 中注册 Webhook 到你的 manager 中。



### 运行和部署controller

安装CRD

```
make install
```

这条命令就是把config/crd里面的资源定义文件apply到我们的集群中去，然后就可以针对我们的集群运行controller了。运行controller的证书和我们连接集群的证书是同一个，因此现在不必担心RBAC权限的问题。

