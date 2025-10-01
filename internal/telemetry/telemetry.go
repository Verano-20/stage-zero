package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/Verano-20/stage-zero/internal/config"
	"github.com/Verano-20/stage-zero/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TelemetryProvider struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	Tracer         trace.Tracer
	Meter          metric.Meter
	AppMetrics     *AppMetrics
	shutdownFuncs  []func(context.Context) error
}

var globalProvider *TelemetryProvider

func InitTelemetry() {
	log := logger.Get()
	log.Info("Initializing Telemetry...")

	globalProvider = &TelemetryProvider{
		shutdownFuncs: make([]func(context.Context) error, 0),
	}

	otelResource, err := newOtelResource()
	if err != nil {
		log.Fatal("Failed to create OpenTelemetry resource", zap.Error(err))
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	var traceExporters []sdktrace.SpanExporter
	var metricReaders []sdkmetric.Reader

	if err := initStdout(&traceExporters, &metricReaders); err != nil {
		log.Fatal("Failed to initialize stdout telemetry", zap.Error(err))
	}
	if err := initOTLP(&traceExporters, &metricReaders); err != nil {
		log.Fatal("Failed to initialize OTLP telemetry", zap.Error(err))
	}

	// TODO: enable no exporters or readers and disable in e2e tests
	if len(traceExporters) == 0 {
		log.Fatal("No trace exporters configured")
	}
	if len(metricReaders) == 0 {
		log.Fatal("No metric readers configured")
	}

	initTracerProvider(otelResource, &traceExporters)
	if err := initMeterProvider(otelResource, &metricReaders); err != nil {
		log.Fatal("Failed to initialize meter provider", zap.Error(err))
	}

	log.Info("Telemetry initialized successfully")
}

func newOtelResource() (*resource.Resource, error) {
	config := config.Get()
	return resource.Merge(resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
			attribute.String("service.instance.id", fmt.Sprintf("%s-%d", config.ServiceName, time.Now().Unix())),
		))
}

func initStdout(traceExporters *[]sdktrace.SpanExporter, metricReaders *[]sdkmetric.Reader) error {
	log := logger.Get()
	config := config.Get()

	if !config.Telemetry.EnableStdout {
		log.Info("Stdout telemetry not enabled, skipping...")
		return nil
	}

	log.Info("Initializing stdout telemetry...")

	stdoutTraceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create stdout trace exporter: %w", err)
	}
	*traceExporters = append(*traceExporters, stdoutTraceExporter)

	stdoutMetricReader, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create stdout metric exporter: %w", err)
	}
	metricReader := sdkmetric.NewPeriodicReader(stdoutMetricReader, sdkmetric.WithInterval(config.Telemetry.MetricInterval))
	*metricReaders = append(*metricReaders, metricReader)

	return nil
}

func initOTLP(traceExporters *[]sdktrace.SpanExporter, metricReaders *[]sdkmetric.Reader) error {
	log := logger.Get()
	config := config.Get()

	if !config.Telemetry.EnableOTLP {
		log.Info("OTLP telemetry not enabled, skipping...")
		return nil
	}

	log.Info("Initializing OTLP telemetry...")

	var otlpTraceOptions []otlptracegrpc.Option
	var otlpMetricOptions []otlpmetricgrpc.Option

	otlpTraceOptions = append(otlpTraceOptions, otlptracegrpc.WithEndpoint(config.Telemetry.OTLPEndpoint))
	otlpMetricOptions = append(otlpMetricOptions, otlpmetricgrpc.WithEndpoint(config.Telemetry.OTLPEndpoint))

	if config.Telemetry.OTLPInsecure {
		otlpTraceOptions = append(otlpTraceOptions, otlptracegrpc.WithInsecure())
		otlpMetricOptions = append(otlpMetricOptions, otlpmetricgrpc.WithInsecure())
	}

	otlpTraceExporter, err := otlptracegrpc.New(context.Background(), otlpTraceOptions...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}
	*traceExporters = append(*traceExporters, otlpTraceExporter)
	globalProvider.shutdownFuncs = append(globalProvider.shutdownFuncs, otlpTraceExporter.Shutdown)

	otlpMetricExporter, err := otlpmetricgrpc.New(context.Background(), otlpMetricOptions...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}
	globalProvider.shutdownFuncs = append(globalProvider.shutdownFuncs, otlpMetricExporter.Shutdown)
	metricReader := sdkmetric.NewPeriodicReader(otlpMetricExporter, sdkmetric.WithInterval(config.Telemetry.MetricInterval))
	*metricReaders = append(*metricReaders, metricReader)

	return nil
}

func initTracerProvider(otelResource *resource.Resource, traceExporters *[]sdktrace.SpanExporter) {
	config := config.Get()

	var traceProcessors []sdktrace.SpanProcessor
	for _, exporter := range *traceExporters {
		processor := sdktrace.NewBatchSpanProcessor(exporter)
		traceProcessors = append(traceProcessors, processor)
		globalProvider.shutdownFuncs = append(globalProvider.shutdownFuncs, processor.Shutdown)
	}

	tracerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(otelResource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // TODO: config for sample rate
	}
	for _, processor := range traceProcessors {
		tracerOptions = append(tracerOptions, sdktrace.WithSpanProcessor(processor))
	}

	globalProvider.TracerProvider = sdktrace.NewTracerProvider(tracerOptions...)
	globalProvider.Tracer = globalProvider.TracerProvider.Tracer(config.ServiceName)
	globalProvider.shutdownFuncs = append(globalProvider.shutdownFuncs, globalProvider.TracerProvider.Shutdown)

	otel.SetTracerProvider(globalProvider.TracerProvider)
}

func initMeterProvider(otelResource *resource.Resource, metricReaders *[]sdkmetric.Reader) error {
	config := config.Get()

	meterOptions := []sdkmetric.Option{
		sdkmetric.WithResource(otelResource),
	}
	for _, reader := range *metricReaders {
		meterOptions = append(meterOptions, sdkmetric.WithReader(reader))
	}

	globalProvider.MeterProvider = sdkmetric.NewMeterProvider(meterOptions...)
	globalProvider.Meter = globalProvider.MeterProvider.Meter(config.ServiceName)
	globalProvider.shutdownFuncs = append(globalProvider.shutdownFuncs, globalProvider.MeterProvider.Shutdown)

	appMetrics, err := NewAppMetrics(globalProvider.Meter)
	if err != nil {
		return fmt.Errorf("failed to create application metrics: %w", err)
	}
	globalProvider.AppMetrics = appMetrics

	otel.SetMeterProvider(globalProvider.MeterProvider)

	return nil
}

func Shutdown() {
	log := logger.Get()
	log.Info("Shutting down Telemetry...")

	if globalProvider == nil {
		log.Warn("Telemetry not initialized, skipping shutdown")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, shutdown := range globalProvider.shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			log.Error("Error during telemetry shutdown", zap.Error(err))
		}
	}

	log.Info("Telemetry shutdown completed")
}

func GetMetrics() *AppMetrics {
	if globalProvider == nil || globalProvider.AppMetrics == nil {
		panic("Metrics not initialized")
	}
	return globalProvider.AppMetrics
}
