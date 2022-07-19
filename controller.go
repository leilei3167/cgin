package main

/*
framework之外的文件均是模拟用户对于框架的使用

*/
import (
	"context"
	"fmt"
	"github.com/leilei3167/cgin/framework"
	"log"
	"time"
)

// FooControllerHandler 调用者视角,如何使用框架
func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan any, 1)

	durationCtx, cancle := context.WithTimeout(c.BaseContext(), time.Second)
	defer cancle()

	go func() {
		//每一个G中都必须有defer recover来捕获异常
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		//do something
		time.Sleep(time.Second * 10)
		c.Json(200, "ok")

		finish <- struct{}{}
	}()

	select {
	//发生panic
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
		//成功处理完毕
	case <-finish:
		fmt.Println("finish")
		//处理时间超时
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}

	return nil
}
