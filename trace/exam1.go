package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
)

func NewTracer(service string) (opentracing.Tracer, io.Closer)  {
	cfg := jaegercfg.Configuration{
		ServiceName: service,
		//采样配置
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 将span发到jaeger-collector的服务地址
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}
	tracer,closer,err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n",err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer,closer
}

/*
1）初始化一个 tracer
2）记录一个简单的 span
3）在span上添加注释信息
*/
func exam1()  {
	// 解析命令行参数
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	// 1. 初始化tracer
	tracer, closer := NewTracer("hello")
	defer closer.Close()

	// 2.开始新的 Span （注意:必须要调用 Finish()方法span才会上传到后端）
	span := tracer.StartSpan("say-hello")
	defer span.Finish()

	helloTo := os.Args[1]
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	// 3.通过tag、log记录注释信息
	// LogFields 和 LogKV底层是调用的同一个方法
	span.SetTag("hello-to", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	span.LogKV("event","println")
	println(helloStr)
}