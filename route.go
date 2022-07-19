package main

import "github.com/leilei3167/cgin/framework"

func registerRouter(core *framework.Core) {
	//需求1,2
	core.GET("/foo", FooControllerHandler)
	//需求3
	leilei := core.Group("/leilei")
	{
		leilei.GET("/age", AgeController)
		//需求4
		leilei.PUT("/:job", AgeController)
	}

}

func AgeController(c *framework.Context) error {
	return nil
}
