package main

import (
	"github.com/leilei3167/cgin/framework"
	"github.com/leilei3167/cgin/framework/middleware"
	"time"
)

// 注册路由规则
func registerRouter(core *framework.Core) {
	// 需求1+2:HTTP方法+静态路由匹配
	core.Use(middleware.Recovery(), middleware.Logger(), middleware.TimoutHandler(time.Second*3))
	core.Get("/user/login", UserLoginController)

	// 需求3:批量通用前缀
	subjectApi := core.Group("/subject")
	subjectApi.Use(middleware.TestGroupMid())
	{

		// 需求4:动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)

		//额外:使需求三可以嵌套
		inner := subjectApi.Group("/inner")
		inner.Get("/find", TestCon)

	}
}

func TestCon(c *framework.Context) error {
	c.Json(200, "ok inner!!!")
	return nil
}
