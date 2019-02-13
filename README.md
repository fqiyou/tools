### tools

#### 项目搭建
```cgo
govendor sync
govendor init

glide init
glide up
glide install

```


#### 使用
1. 监控export
    1. raid 磁盘监控
    ```cgo
    # raid 磁盘监控,打包
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/run/monitor/raid_disk foo/monitor/host/raid_disk/main.go
    
    ```
    ```cgo
    [root@bj-sjs-009 raid]# ./raid_disk --help
    Usage of ./raid_disk:
      -metric.namespace string
            Prometheus metrics namespace, as the prefix of metrics name (default "sjs")
      -web.listen-port string
            An port to listen on for web interface and telemetry. (default "9001")
      -web.telemetry-path string
            A path under which to expose metrics. (default "/metrics")
    [root@bj-sjs-009 raid]# /opt/soft/monitor/raid/raid_disk   --metric.namespace sjs --web.listen-port 9992 --web.telemetry-path /metrics
    INFO[0000] Starting Server at http://localhost:9992/metrics  source="raid_disk/main.go:41"

    ```

