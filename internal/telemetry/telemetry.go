package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/Verano-20/go-crud/internal/config"
	"github.com/Verano-20/go-crud/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TelemetryProvider struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	Tracer         trace.Tracer
	Meter          metric.Meter
	shutdownFuncs  []func(context.Context) error
}

var globalProvider *TelemetryProvider

func InitTelemetry() {
	config := config.Get()
	log := logger.Get()

	log.Info("Initializing Telemetry...")

	resource, err := newResource(config)
	if err != nil {
		log.Fatal("Failed to create telemetry resource", zap.Error(err))
	}

	tp := &TelemetryProvider{
		shutdownFuncs: make([]func(context.Context) error, 0),
	}

	if err := tp.setupTracing(config, resource); err != nil {
		log.Fatal("Failed to setup tracing", zap.Error(err))
	}

	if err := tp.setupMetrics(config, resource); err != nil {
		log.Fatal("Failed to setup metrics", zap.Error(err))
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	globalProvider = tp
	log.Info("Telemetry initialized successfully")
}

func (tp *TelemetryProvider) setupTracing(config *config.Config, resource *resource.Resource) error {
	var exporters []sdktrace.SpanExporter

	// Stdout exporter for development
	if config.Telemetry.EnableStdoutTrace {
		stdoutExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}
		exporters = append(exporters, stdoutExporter)
	}

	// OTLP exporter for production/cloud
	if config.Telemetry.EnableOTLP {
		var otlpOptions []otlptracehttp.Option
		otlpOptions = append(otlpOptions, otlptracehttp.WithEndpoint(config.Telemetry.OTLPEndpoint))

		if config.Telemetry.OTLPInsecure {
			otlpOptions = append(otlpOptions, otlptracehttp.WithInsecure())
		}

		otlpExporter, err := otlptracehttp.New(context.Background(), otlpOptions...)
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
		exporters = append(exporters, otlpExporter)
		tp.shutdownFuncs = append(tp.shutdownFuncs, otlpExporter.Shutdown)
	}

	var processors []sdktrace.SpanProcessor
	for _, exporter := range exporters {
		processor := sdktrace.NewBatchSpanProcessor(exporter)
		processors = append(processors, processor)
		tp.shutdownFuncs = append(tp.shutdownFuncs, processor.Shutdown)
	}

	tracerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	}
	for _, processor := range processors {
		tracerOptions = append(tracerOptions, sdktrace.WithSpanProcessor(processor))
	}

	tp.TracerProvider = sdktrace.NewTracerProvider(tracerOptions...)

	// Set global tracer provider
	otel.SetTracerProvider(tp.TracerProvider)

	tp.Tracer = tp.TracerProvider.Tracer(config.ServiceName)
	tp.shutdownFuncs = append(tp.shutdownFuncs, tp.TracerProvider.Shutdown)

	return nil
}

// TODO: Add actual metrics
func (tp *TelemetryProvider) setupMetrics(config *config.Config, resource *resource.Resource) error {
	var readers []sdkmetric.Reader

	// Stdout exporter for development
	if config.Telemetry.EnableStdoutTrace { // Use same flag for simplicity
		stdoutExporter, err := stdoutmetric.New()
		if err != nil {
			return fmt.Errorf("failed to create stdout metric exporter: %w", err)
		}
		reader := sdkmetric.NewPeriodicReader(stdoutExporter, sdkmetric.WithInterval(30*time.Second)) // TODO: config
		readers = append(readers, reader)
	}

	// OTLP exporter for production/cloud
	if config.Telemetry.EnableOTLP {
		var otlpOptions []otlpmetrichttp.Option
		otlpOptions = append(otlpOptions, otlpmetrichttp.WithEndpoint(config.Telemetry.OTLPEndpoint))

		if config.Telemetry.OTLPInsecure {
			otlpOptions = append(otlpOptions, otlpmetrichttp.WithInsecure())
		}

		otlpExporter, err := otlpmetrichttp.New(context.Background(), otlpOptions...)
		if err != nil {
			return fmt.Errorf("failed to create OTLP metric exporter: %w", err)
		}

		reader := sdkmetric.NewPeriodicReader(otlpExporter, sdkmetric.WithInterval(30*time.Second)) // TODO: config
		readers = append(readers, reader)
		tp.shutdownFuncs = append(tp.shutdownFuncs, otlpExporter.Shutdown)
	}

	meterOptions := []sdkmetric.Option{
		sdkmetric.WithResource(resource),
	}
	for _, reader := range readers {
		meterOptions = append(meterOptions, sdkmetric.WithReader(reader))
	}

	tp.MeterProvider = sdkmetric.NewMeterProvider(meterOptions...)

	// Set global meter provider
	otel.SetMeterProvider(tp.MeterProvider)

	tp.Meter = tp.MeterProvider.Meter(config.ServiceName)
	tp.shutdownFuncs = append(tp.shutdownFuncs, tp.MeterProvider.Shutdown)

	return nil
}

func (tp *TelemetryProvider) Shutdown(ctx context.Context) error {
	log := logger.Get()
	log.Info("Shutting down Telemetry...")

	for _, shutdown := range tp.shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			log.Error("Error during telemetry shutdown", zap.Error(err))
		}
	}

	log.Info("Telemetry shutdown completed")
	return nil
}

func newResource(config *config.Config) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
			attribute.String("service.instance.id", fmt.Sprintf("%s-%d", config.ServiceName, time.Now().Unix())),
		))
}

func Shutdown() error {
	if globalProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return globalProvider.Shutdown(ctx)
	}
	return nil
}
