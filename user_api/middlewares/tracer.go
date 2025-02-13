package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TracerMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan(c.Request.Method + "-" + c.Request.URL.Path)
		defer span.Finish()

		c.Set("span", span)
		c.Next()
	}
}

func GrpcTracerInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		parentSpan := opentracing.SpanFromContext(ctx)
		var span opentracing.Span
		if parentSpan == nil {
			zap.S().Infof(method + " with no parent span")
			span = tracer.StartSpan(method)
		} else {
			span = tracer.StartSpan(
				method,
				opentracing.ChildOf(parentSpan.Context()),
			)
		}
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
