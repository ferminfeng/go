package main

import (
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func main() {

	err := profiler.Start(
		profiler.WithService("fermin-service-profiler"),
		profiler.WithEnv("test"),
		profiler.WithVersion("dd-1.0.0"),
		profiler.WithTags("k:1", "k:2"),
		profiler.WithAgentAddr("localhost:9529"), // DataKit url
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			// The profiles below are disabled by default to keep overhead
			// low, but can be enabled as needed.

			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	tracer.Start(
		tracer.WithEnv("test"),
		tracer.WithService("fermin-service-tracer"),
		tracer.WithServiceVersion("1.2.3"),
		tracer.WithGlobalTag("project", "fermin- add-ddtrace-in-golang-project"),
		tracer.WithAgentAddr("localhost:9529"), // DataKit url
	)

	tick := time.NewTicker(time.Second)

	defer func() {
		profiler.Stop()
		tracer.Stop()
		tick.Stop()
	}()

	// your-app-main-entry...
	//for {
	runApp()
	runAppWithError()

	//	select {
	//	case <-tick.C:
	//	}
	//}

	// 开启 mutex 和 block 性能采集

	// 设置采集频率，即 1/rate 的事件被采集， 如设置为 0 或小于 0 的数值，是不进行采集的
	var rate = 1

	// enable mutex profiling
	runtime.SetMutexProfileFraction(rate)

	// enable block profiling
	runtime.SetBlockProfileRate(rate)
	http.ListenAndServe(":6060", nil)

}

func runApp() {
	var err error
	// Start a root span.
	span := tracer.StartSpan("get.data")
	defer span.Finish(tracer.WithError(err))

	// Create a child of it, computing the time needed to read a file.
	child := tracer.StartSpan("read.file", tracer.ChildOf(span.Context()))
	child.SetTag(ext.ResourceName, os.Args[0])

	// Perform an operation.
	var bts []byte
	bts, err = ioutil.ReadFile(os.Args[0])
	span.SetTag("file_len", len(bts))
	child.Finish(tracer.WithError(err))
}

func runAppWithError() {
	var err error
	// Start a root span.
	span := tracer.StartSpan("get.data")

	// Create a child of it, computing the time needed to read a file.
	child := tracer.StartSpan("read.file", tracer.ChildOf(span.Context()))
	child.SetTag(ext.ResourceName, "somefile-not-found.go")

	defer func() {
		child.Finish(tracer.WithError(err))
		span.Finish(tracer.WithError(err))
	}()

	// Perform an error operation.
	if _, err = ioutil.ReadFile("somefile-not-found.go"); err != nil {
		// error handle
	}
}
