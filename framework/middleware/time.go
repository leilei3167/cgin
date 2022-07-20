package middleware

import (
	"github.com/leilei3167/cgin/framework"
	"log"
	"time"
)

func Logger() framework.ControllerHandler {
	return func(c *framework.Context) error {
		start := time.Now()

		c.Next()

		log.Printf("uri:%s spent:%v", c.GetRequest().URL.Path, time.Since(start))

		return nil
	}
}
