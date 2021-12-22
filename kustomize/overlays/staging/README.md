演示环境，定义新的名称前缀以及一些不同的标签

```
xing@xingdeMacBook-Pro kustomize % kustomize build overlays/staging 
apiVersion: v1
data:
  altGreeting: Have a pineapple!
  enableRisky: "true"
kind: ConfigMap
metadata:
  annotations:
    note: Hello, I am staging!
  labels:
    app: hello
    org: acmeCorporation
    variant: staging
  name: staging-the-map
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: Hello, I am staging!
  labels:
    app: hello
    org: acmeCorporation
    variant: staging
  name: staging-the-service
spec:
  ports:
  - port: 8666
    protocol: TCP
    targetPort: 8080
  selector:
    app: hello
    deployment: hello
    org: acmeCorporation
    variant: staging
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    note: Hello, I am staging!
  labels:
    app: hello
    org: acmeCorporation
    variant: staging
  name: staging-the-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hello
      deployment: hello
      org: acmeCorporation
      variant: staging
  template:
    metadata:
      annotations:
        note: Hello, I am staging!
      labels:
        app: hello
        deployment: hello
        org: acmeCorporation
        variant: staging
    spec:
      containers:
      - command:
        - /hello
        - --port=8080
        - --enableRiskyFeature=$(ENABLE_RISKY)
        env:
        - name: ALT_GREETING
          valueFrom:
            configMapKeyRef:
              key: altGreeting
              name: staging-the-map
        - name: ENABLE_RISKY
          valueFrom:
            configMapKeyRef:
              key: enableRisky
              name: staging-the-map
        image: monopole/hello:1
        name: the-container
        ports:
        - containerPort: 8080
```

添加一个 ConfigMap 自定义把 base 中的 ConfigMap 中的 "*Good Morning!*" 变成 "Good Night!"