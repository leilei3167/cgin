package middleware

import (
	"github.com/leilei3167/cgin/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		c.Next() //在执行后续的处理链出现panic将被捕获
		return nil
	}
}
