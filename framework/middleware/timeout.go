package middleware

//# v1版本,利用函数嵌套的形式,不符合要求
/*
func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	//接收业务逻辑作为参数,利用函数回调 执行,返回值是匿名函数
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan any, 1)

		durationCtx, cancle := context.WithTimeout(c.BaseContext(), d)
		defer cancle()

		go func() {
			//每一个G中都必须有defer recover来捕获异常
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			//此处执行业务逻辑
			fun(c)

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
}
*/

//# v2使用pipeline的思想
