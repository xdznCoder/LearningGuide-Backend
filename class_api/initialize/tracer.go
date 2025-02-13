package initialize

import (
	"LearningGuide/user_api/global"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"io"
)

func InitTracer() io.Closer {
	sampler := jaeger.NewConstSampler(true)
	sender := transport.NewHTTPTransport(fmt.Sprintf("http://%s:%d/api/traces",
		global.ServerConfig.Jaeger.Host,
		global.ServerConfig.Jaeger.Port,
	))

	reporter := jaeger.NewRemoteReporter(sender)

	tracer, closer := jaeger.NewTracer(global.ServerConfig.Name, sampler, reporter)

	opentracing.SetGlobalTracer(tracer)

	return closer
}
