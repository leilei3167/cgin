package main

import "github.com/leilei3167/cgin/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
