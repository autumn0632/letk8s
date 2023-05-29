# operator 简介

添加自定义资源步骤：
1. 编写crd
2. 开发crd对应的控制器



# kubebuilder 原理

Kubebuilder 是一个使用CRDs构建 K8s API 的 SDK，主要是：

* 提供脚手架工具初始化 CRDs 工程，自动生成 boilerplate 代码和配置；
* 提供代码库封装底层的 K8s go-client；

方便用户从零开始开发 CRDs，Controllers 和 Admission Webhooks 来扩展 K8s。

## 核心概念



### GVKs&GVRS

**GVK** = GroupVersionKind
**GVR** = GroupVersionResource。

API Group & Versions（GV）
API Group 是相关 API 功能的集合，每个 Group 拥有一或多个 Versions，用于接口的演进。

**Kinds & Resources**

* 每个 GV 都包含多个 API 类型，称为 Kinds，在不同的 Versions 之间同一个 Kind 定义可能不同
* Resource 是 Kind 的对象标识（resource type），一般来说 Kinds 和 Resources 是 1:1 的，比如 pods Resource 对应 Pod Kind，但是有时候相同的 Kind 可能对应多个 Resources，比如 Scale Kind 可能对应很多 Resources：deployments/scale，replicasets/scale，对于 CRD 来说，只会是 1:1 的关系。

每一个 GVK 都关联着一个 package 中给定的 root Go type，比如 apps/v1/Deployment 就关联着 K8s 源码里面 k8s.io/api/apps/v1 package 中的 Deployment struct，我们提交的各类资源定义 YAML 文件都需要写：

* apiVersion：这个就是 GV 。
* kind：这个就是 K。

**根据 GVK K8s 就能找到你到底要创建什么类型的资源，根据你定义的 Spec 创建好资源之后就成为了 Resource，也就是 GVR。**GVK/GVR 就是 K8s 资源的坐标，是我们创建/删除/修改/读取资源的基础。

**Scheme**

每一组 Controllers 都需要一个 Scheme，提供了 Kinds 与对应 Go types 的映射，也就是说给定 Go type 就知道他的 GVK，给定 GVK 就知道他的 Go type。

### 代码组件

**Manager**

Kubebuilder 的核心组件，具有 3 个职责：

* 负责运行所有的 Controllers；
* 初始化共享 caches，包含 listAndWatch 功能；
* 初始化 clients 用于与 Api Server 通信。

**Cache**

Kubebuilder 的核心组件，负责在 Controller 进程里面根据 Scheme 同步 Api Server 中所有该 Controller 关心 GVKs 的 GVRs，其核心是 GVK -> Informer 的映射，**Informer 会负责监听对应 GVK 的 GVRs 的创建/删除/更新操作，以触发 Controller 的 Reconcile 逻辑。**

**Controller**

Kubebuidler 为我们生成的脚手架文件，我们只需要实现 Reconcile 方法即可。（实现自己的逻辑）

**Clients**

在实现 Controller 的时候不可避免地需要对某些资源类型进行创建/删除/更新，就是通过该 Clients 实现的，其中查询功能实际查询是本地的 Cache，写操作直接访问 Api Server。

**Index**

由于 Controller 经常要对 Cache 进行查询，Kubebuilder 提供 Index utility 给 Cache 加索引提升查询效率。

**OwnerReference**

K8s GC 在删除一个对象时，任何 ownerReference 是该对象的对象都会被清除，与此同时，Kubebuidler 支持所有对象的变更都会触发 Owner 对象 controller 的 Reconcile 方法。


# 实战

1. 初始化脚手架框架

    > kubebuilder init --domain autumn.io --repo github.com/autumn0632/smartddi-crd

2. 创建api

    > kubebuilder create api --group test --version v1 --kind  SmartDDI

3. 定义CRD

    修改Status和Spec定义
    ```go
    // api/v1/smartddi_type.go

    // SmartDDISpec defines the desired state of SmartDDI
    type SmartDDISpec struct {
	    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	    // Important: Run "make" to regenerate code after modifying this file

	    // Foo is an example field of SmartDDI. Edit smartddi_types.go to remove/update
	    // Foo string `json:"foo,omitempty"`   
        Version stirng `json:"version"`
    }

    // SmartDDIStatus defines the observed state of SmartDDI
    type SmartDDIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
        Version stirng `json:"version"`
    }

    ```

    4. 编写Controller逻辑

    在文件`controllers/smartddi_controller.go`中，在`Reconcile`函数里添加逻辑

    5. 生成并部署相关资源文件
    
    > make install

    6. make run 运行

    7. 部署cr资源，查看日志


# 创建自定义资源并应用到k8s
让 K8s 知道有这个资源及其结构属性，在用户提交该自定义资源的定义时（通常是 YAML 文件定义），K8s 能够成功校验该资源并创建出对应的`Go struct`进行持久化，同时触发控制器的调谐逻辑。



# 编写控制器监控自定义资源
