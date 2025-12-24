package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "net/http/pprof"
)

func main() {

	shutdown := initOpenTelemetry()
	defer func() {
		shutdown()
	}()

	log.Println("connect ...")

	http.HandleFunc("/user", web)

	err := http.ListenAndServe(":14317", nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("服务启动失败: %v", err))
	}

	time.Sleep(time.Minute * 2)
	os.Exit(0)
}

// 初始化链路追踪
func initOpenTelemetry() func() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("fermin-OpenTelemetry"),
			semconv.DeploymentEnvironmentKey.String("test"),
		),
	)

	if err != nil {
		log.Fatal("【dataKit】初始化失败")
	}

	var bsp sdktrace.SpanProcessor

	conn, err := grpc.DialContext(ctx, "127.0.0.1:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatal("【dataKit】failed to create gRPC connection to collector")
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))

	if err != nil {
		log.Fatal("【dataKit】failed to create trace exporter")
	}

	bsp = sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		err = tracerProvider.Shutdown(ctx)
		if err != nil {
			log.Fatal("【dataKit】failed to shutdown TracerProvider")
		}
	}
}

var tracer = otel.Tracer("tracer_user_login")

// web handler 处理请求数据
func web(w http.ResponseWriter, r *http.Request) {
	_, spanA := InitBuryingPoint("demo-error-first")
	spanA.SetStatus(codes.Error, "span-first is error")
	spanA.End()

	ctx, span := InitBuryingPoint("demo-first")
	span.SetStatus(codes.Error, "span-first is ok")

	defer span.End()

	// 调用服务层
	service(ctx)

	testPoint(ctx, "aaaa")
	testPoint(ctx, "bbbb")

	_, _ = w.Write([]byte("ok"))
	// ... 接收客户端请求
	log.Println("finish ...")
}

func testPoint(ctx context.Context, name string) {
	_, span := AddBuryingPoint(ctx, name)
	defer span.End()
}

// service 调用 service 层处理业务
func service(ctx context.Context) {

	ctx, iSpan := AddBuryingPoint(ctx, "demo-service")
	defer iSpan.End()

	//iSpan.SetStatus(codes.Ok, "demo-service is ok")

	testPoint(ctx, "cccc")
	dao(ctx)

	testPoint(ctx, "dddd")
}

// dao 数据访问层
func dao(ctx context.Context) {
	testPoint(ctx, "eeee")

	ctx, iSpan := AddBuryingPoint(ctx, "demo-dao")
	//iSpan.SetStatus(codes.Ok, "demo-dao is ok")

	testPoint(ctx, "ffff")

	// 创建子 span 查询数据库等操作
	_, sqlSpan := AddBuryingPoint(ctx, "demo-sql")
	// 成功时返回标识
	//sqlSpan.SetStatus(codes.Ok, "demo-sql is ok")

	testPoint(ctx, "gggg")

	// 失败时返回标识
	//sqlSpan.SetStatus(codes.Error, "demo-sql 失败了，来看看失败原因")
	FinishBuryingPoint(sqlSpan, false)

	defer func() {
		//sqlSpan.End()

		iSpan.End()
	}()
}

// InitBuryingPoint 开启埋点
func InitBuryingPoint(spanName string) (context.Context, trace.Span) {
	commonLabels := []attribute.KeyValue{attribute.String("key1", "自定义键值对-val")}

	ctx := context.Background()

	ctx, span := tracer.Start(
		ctx,
		spanName,
		trace.WithAttributes(commonLabels...),
	)
	return ctx, span
}

// AddBuryingPoint 埋点
func AddBuryingPoint(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}

// FinishBuryingPoint 结束埋点
func FinishBuryingPoint(dataKitSpan trace.Span, res bool) {

	dataKitSpan.End()

	dataKitSpan.SetAttributes(attribute.String("http_response_content", "哈哈哈哈哈"))
	if res {
		dataKitSpan.SetStatus(codes.Ok, "is ok")
	} else {
		dataKitSpan.SetStatus(codes.Error, "is fail")
	}
}
