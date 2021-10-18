# apollo-tools

Apollo-tools 是一个命令行的小工具，用于获取 Apollo 配置中心的配置，输出到指定文件。

一次性执行，输出应用配置

```bash
apollo-tools -app-id app -server-addr http://127.0.0.1:8088 -format "%s: %s" -output config.yaml
```

持续运行，当配置变更时，更新应用

```bash
apollo-tools -app-id app -server-addr http://127.0.0.1:8088 -format "%s=%s" -forever -on-change 'cat .env'  -output .env
```

在 Kubernetes 中为不支持 Apollo 配置中心的项目生成配置文件

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: APP NAME
  labels:
    app: APP NAME
spec:
  selector:
    matchLabels:
      app: APP NAME
  replicas: 1
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: APP NAME
    spec:
      initContainers:
        - name: apollo-config-loader
          image: mylxsw/apollo-tools:1.0
          args:
            - "-server-addr"
            - "http://apollo-config:8080"
            - "-app-id"
            - "APP ID"
            - "-format"
            - "%s: %s"
            - "-output"
            - "/data/webroot/config.yaml"
          volumeMounts:
            - name: config-dir
              mountPath: /data/webroot
      containers:
        - name: web-tools
          image: APP IMAGE:VERSION
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          volumeMounts:
            - name: config-dir
              mountPath: /data/webroot
      volumes:
        - name: config-dir
          emptyDir: {}
```
