package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"os"
)


//通过opentracing.ChildOf(rootSpan.Context())保留span之间的因果关系。
func exam2()  {
	// 解析命令行参数
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	// 1.初始化 tracer
	tracer, closer := NewTracer("hello")
	defer closer.Close()
	// 2.开始新的 Span （注意:必须要调用 Finish()方法span才会上传到后端）
	span := tracer.StartSpan("say-hello")
	defer span.Finish()

	helloTo := os.Args[1]
	helloStr := formatString(span, helloTo)
	printHello(span, helloStr)
}

func formatString(span opentracing.Span, helloTo string) string  {
	//记录调用关系,就是构造span的过程，然后把构造的span通过finish()传到collector里
	childSpan := span.Tracer().StartSpan(
		"formatString",
		opentracing.ChildOf(span.Context()),
		)
	defer childSpan.Finish()

	//操作逻辑
	return fmt.Sprintf("Hello, %s!",helloTo)
}

func printHello(span opentracing.Span, helloStr string)  {
	//记录调用关系
	childSpan := span.Tracer().StartSpan(
		"printHello",
		opentracing.ChildOf(span.Context()),
		)
	defer childSpan.Finish()

	//操作逻辑
	println(helloStr)
}