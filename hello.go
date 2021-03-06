package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"time"
)

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {

	localAgentHostPort := os.Getenv("LOCAL_AGENT_HOST_PORT")
	collectorHostPort := os.Getenv("COLLECTOR_HOST_PORT")

	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: localAgentHostPort,
			CollectorEndpoint:  collectorHostPort,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}

func main() {

	for {

		tracer, closer := initJaeger("hello-world")

		opentracing.SetGlobalTracer(tracer)

		span := tracer.StartSpan("say-hello")

		ctx := opentracing.ContextWithSpan(context.Background(), span)

		helloStr := formatString(ctx, "Ehsan")

		span.SetTag("hello-to", "Ehsan")

		printHello(ctx, helloStr)

		span.Finish()
		_ = closer.Close()

		time.Sleep(10000 * time.Millisecond)
	}
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}
