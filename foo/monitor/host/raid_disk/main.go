
package main

import (
	"flag"
	"github.com/fqiyou/tools/foo/monitor/host/raid_disk/collector"
	"github.com/fqiyou/tools/foo/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// 命令行参数
	listenAddr  = flag.String("web.listen-port", "9001", "An port to listen on for web interface and telemetry.")
	metricsPath = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace = flag.String("metric.namespace", "sjs", "Prometheus metrics namespace, as the prefix of metrics name")
)

func main() {
	flag.Parse()

	metrics := collector.NewMetrics(*metricsNamespace)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	util.Log.Info("Starting Server at http://localhost:", *listenAddr, *metricsPath)
	util.Log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
}