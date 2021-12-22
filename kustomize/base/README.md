run `kustomize build` 会在标准输出中为每个资源对象贴上 `app: hello` 的标签
可以运行`kustomize build base | kubectl apply -f -` 命令直接部署在集群中

