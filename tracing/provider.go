package tracing

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// NewTracerProvider 创建一个链路追踪器
func NewTracerProvider(exporterName, endpoint, serviceName, instanceId, version string, sampler float64) *trace.TracerProvider {
	if instanceId == "" {
		ud, _ := uuid.NewUUID()
		instanceId = ud.String()
	}
	if version == "" {
		version = "x.x.x"
	}

	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(sampler))),
		trace.WithResource(resource.NewSchemaless(
			semConv.ServiceNameKey.String(serviceName),
			semConv.ServiceInstanceIDKey.String(instanceId),
			semConv.ServiceVersionKey.String(version),
		)),
	}

	if len(endpoint) > 0 {
		exp, err := NewExporter(exporterName, endpoint, true)
		if err != nil {
			panic(err)
		}

		opts = append(opts, trace.WithBatcher(exp))
	}

	return trace.NewTracerProvider(opts...)
}
