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
	"time"
)

/*
追踪 gRPC 则通过拦截器实现。
这里使用使用 gRPC 的metadata 来做载体。
*/

// ClientInterceptor grpc client
func ClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor  {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 一个RPC调用的服务端的span和客户端的span构成ChildOf关系
		var parentCtx opentracing.SpanContext
		parentSpan := opentracing.SpanFromContext(ctx)
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}
		span := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			opentracing.Tag{Key: string(ext.Component),Value: "gRPC Client"},
			ext.SpanKindRPCClient,
			)

		//简单写法
		//span,_ := opentracing.StartSpanFromContext(ctx,method,
		//	opentracing.Tag{Key: string(ext.Component),Value: "gRPC"},
		//	ext.SpanKindRPCClient,
		//	)

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
			span.LogFields(log.String("inject span error: %v",err.Error()))
		}

		//调用远程方法
		newCtx := metadata.NewOutgoingContext(ctx,md)
		err = invoker(newCtx, method, req,reply,cc,opts...)
		if err != nil {
			log.Error(err)
		}
		return err
	}

}

func ClientMain(serviceName,serverAddr string){
	tracer,closer := NewTracer(serviceName)
	defer closer.Close()

	ctx,_ := context.WithTimeout(context.Background(),5*time.Second)
	conn,err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ClientInterceptor(tracer)),
		)
	if err != nil {
		log.Error(err)
	}

	defer conn.Close()

	//.....gRPC客户端逻辑
}

func ServerMain(serviceName string)  {
	tracer,closer := NewTracer(serviceName)
	defer closer.Close()

	s := grpc.NewServer(grpc.UnaryInterceptor(ServerInterceptor(tracer)))

	// server处理逻辑
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
				opentracing.Tag{Key: string(ext.Component),Value: "gRPC Server"},
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

/*
官方制定的carrier有两种，TextMapCarrier和HTTPHeadersCarrier
// TextMapCarrier allows the use of regular map[string]string
// as both TextMapWriter and TextMapReader.
type TextMapCarrier map[string]string

// HTTPHeadersCarrier satisfies both TextMapWriter and TextMapReader.
//
// Example usage for server side:
//
//     carrier := opentracing.HTTPHeadersCarrier(httpReq.Header)
//     clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
//
// Example usage for client side:
//
//     carrier := opentracing.HTTPHeadersCarrier(httpReq.Header)
//     err := tracer.Inject(
//         span.Context(),
//         opentracing.HTTPHeaders,
//         carrier)
//
type HTTPHeadersCarrier http.Header
*/

