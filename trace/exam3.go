package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"os"
)


//前面虽然保留的 span 的因果关系，但是需要在各个方法中传递 span。
//这可能会污染整个程序，我们可以借助 Go 语言中的 context.Context对象来进行传递。
func exam3()  {
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

	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx,span)

	helloStr := formatString2(ctx, helloTo)
	printHello2(ctx, helloStr)
}

func formatString2(ctx context.Context, helloTo string) string  {
	/*
	opentracing.StartSpanFromContext()返回的第二个参数是子ctx
	如果需要的话可以将该子ctx继续往下传递，而不是传递父ctx。
	*/
	//这里拿到的就是子span
	span, _ := opentracing.StartSpanFromContext(ctx,"formatString2")

	defer span.Finish()

	//操作逻辑
	return fmt.Sprintf("Hello, %s!",helloTo)
}

func printHello2(ctx context.Context, helloStr string)  {
	//记录调用关系
	span, _ := opentracing.StartSpanFromContext(ctx,"printHello2")

	defer span.Finish()

	//操作逻辑
	println(helloStr)
}