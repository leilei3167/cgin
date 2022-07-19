package framework

/*
在Server结构中 包含了BaseContext和ConnContext 两个函数类型,是两个Context 的修改点
第一个默认是BackGround,用于指定最基本的Context 的类型,Conn则是对每个请求的Context的设置(设置后加入到conn中保存)


因为http调用是多阶段的调用,为了避免某个调用超时等待导致服务崩溃,必须设置相应的超时取消设置,Context就是最佳的选择,02节目标就是实现一个自定义的Context功能
不是想直接使用标准库的 Context，因为它完全是标准库 Context 接口的实现，只能控制链条结束，封装性并不够

在框架里，我们需要有更强大的 Context，除了可以控制超时之外，常用的功能比如获取请求、返回结果、实现标准库的 Context 接口，也都要有。

功能一:处理请求返回结果
用标准库操作http请求主要是调用了 http.Request 和 http.ResponseWriter ，实现 WebService 接收和处理协议文本的功能。
但这两个结构提供的接口粒度太细了，需要使用者非常熟悉这两个结构的内部字段，比如 response 里设置 Header 和设置 Body 的函数，
用起来肯定体验不好。

如果能够将这两个结构的方法 或者函数进一步封装,对外暴露语义化高的接口函数,整个框架的易用性将大大提升!

因此在context.go中构建结构体 封装req和resp


功能二:实现Context接口,实现对处理器的控制


自己封装的 Context 最终需要提供四类功能函数：
base 封装基本的函数功能，比如获取 http.Request 结构
context 实现标准 Context 接口
request 封装了 http.Request 的对外接口
response 封装了 http.ResponseWriter 对外接口

*/

/*
调用者视角:

为单个请求设置超时:
	如何使用自定义 Context 设置超时呢？结合前面分析的标准库思路，我们三步走完成：
1.继承 request 的 Context，创建出一个设置超时时间的 Context；
2.创建一个新的 Goroutine 来处理具体的业务逻辑；
3.设计事件处理顺序，当前 Goroutine 监听超时时间 Contex 的 Done() 事件，和具体的业务处理结束事件，哪个先到就先处理哪个。

接下来就在业务的controller.go中来进行创建

关键点:
	-父G无法铺货子G的panic,每一个G创建时,在其内部都要用defer -recover进行异常捕获,否则单个协程异常可能导致整个系统崩溃!
	-异常,超时触发时,需要写入response,保证并发安全(锁)
	-超时触发后,已经向resp中写入过数据了,其他G也要操作的话,是否出现重复写入(超时标记变量)?



*/
