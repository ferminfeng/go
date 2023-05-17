package main

import (
	"context"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	shutdown := initProvider()
	defer func() {
		shutdown()
		profiler.Stop()
	}()

	log.Println("connect ...")

	// 性能采集
	profile()

	http.HandleFunc("/user", web)
	handleErr(http.ListenAndServe(":14317", nil), "open server")

	//handleErr(http.ListenAndServe(":6060", nil), "性能采集")

	time.Sleep(time.Minute * 2)
	os.Exit(0)
}

func profile() {

	log.Println("性能采集 start ...")

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

	handleErr(err, "性能采集")

	// 开启 mutex 和 block 性能采集 设置采集频率，即 1/rate 的事件被采集， 如设置为 0 或小于 0 的数值，是不进行采集的
	var rate = 1
	runtime.SetMutexProfileFraction(rate)
	runtime.SetBlockProfileRate(rate)
	log.Println("性能采集 end ...")
}

func initProvider() func() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("fermin-OpenTelemetry"),
			// semconv.FaaSIDKey.String(""),
		),
		//resource.WithOS(), // and so on ...
	)

	handleErr(err, "failed to create resource")
	var bsp sdktrace.SpanProcessor

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	conn, err := grpc.DialContext(ctx, "127.0.0.1:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	handleErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	handleErr(err, "failed to create trace exporter")

	bsp = sdktrace.NewBatchSpanProcessor(traceExporter)

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")

		time.Sleep(time.Second)
	}
}

var tracer = otel.Tracer("tracer_user_login")

// web handler 处理请求数据
func web(w http.ResponseWriter, r *http.Request) {
	// ... 接收客户端请求
	log.Println("doing web")

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []attribute.KeyValue{attribute.String("key1", "val1")}

	// work begins
	ctx, span := tracer.Start(
		context.Background(),
		"span-Example",
		trace.WithAttributes(commonLabels...),
	)

	defer span.End()
	<-time.After(time.Millisecond * 50)
	service(ctx)

	log.Printf("Doing really hard work")
	<-time.After(time.Millisecond * 40)

	log.Printf("Done!")
	w.Write([]byte("ok"))
}

// service 调用 service 层处理业务
func service(ctx context.Context) {
	log.Println("service")

	ctx1, iSpan := tracer.Start(ctx, "Sample-service")

	<-time.After(time.Second / 2) // do something...

	dao(ctx1)

	iSpan.End()
}

// dao 数据访问层
func dao(ctx context.Context) {
	log.Println("dao")
	ctxD, iSpan := tracer.Start(ctx, "Sample-dao")
	<-time.After(time.Second / 2)

	// 创建子 span 查询数据库等操作
	_, sqlSpan := tracer.Start(ctxD, "do_sql")
	sqlSpan.SetStatus(codes.Ok, "is ok") //
	<-time.After(time.Second)

	sqlSpan.End()

	iSpan.End()
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
