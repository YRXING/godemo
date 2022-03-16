package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

/*
追踪 gRPC 则通过拦截器实现。
这里使用使用 gRPC 的metadata 来做载体。
*/

// ClientInterceptor grpc client
func ClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor  {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		span,_ := opentracing.StartSpanFromContext(ctx,"call gRPC",
			opentracing.Tag{Key: string(ext.Component),Value: "gRPC"},
			ext.SpanKindRPCClient,
			)
		defer span.Finish()

		md,ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}else {
			md = md.Copy()
		}

		//把span信息注册到MD这个carrier中去
		err := tracer.Inject(span.Context(),opentracing.TextMap,MDReaderWriter{md})
		if err != nil {
			span.LogFields(log.String("call-error",err.Error()))
		}
		return err
	}

}

func ServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor  {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
		md,ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		//服务端拦截则是把在MD中的span提取出来
		spanContext,err := tracer.Extract(opentracing.TextMap,MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			fmt.Println("extract from metadata error: ",err)
		}else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component),Value: "gRPC"},
				ext.SpanKindRPCServer,
				)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx,span)
		}

		return handler(ctx,req)
	}
}

//为了做载体，必须要实现 opentracing.TextMapWriter opentracing.TextMapReader 这两个接口
//For Tracer.Inject(): the carrier must be a `TextMapWriter`.
//For Tracer.Extract(): the carrier must be a `TextMapReader`.
type MDReaderWriter struct {
	metadata.MD
}

func (m MDReaderWriter) ForeachKey(handler func(key, val string) error) error  {
	for k,vs := range m.MD {
		for _,v := range vs{
			if err := handler(k,v);err!= nil {
				return err
			}
		}
	}
	return nil
}

func (m MDReaderWriter) Set(key,val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key],val)
}



