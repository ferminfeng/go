package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
	_ "net/http/pprof" // 自动注册路由到 http.DefaultServeMux
)

func main() {
	// 暴露 /metrics 接口，Prometheus 会自动抓取 runtime 指标
	// 包括：go_goroutines, go_memstats_alloc_bytes, process_cpu_seconds_total 等
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	fmt.Print("启动成功。。。")
	// 你的业务代码...

	for i := 0; i < 100; i++ {
		go func(i int) {
			for {
				select {
				default:
					fmt.Println(fmt.Sprintf("i[%d], time:%d", i, time.Now().Unix()))

					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}

	select {}
}
