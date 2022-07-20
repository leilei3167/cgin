package main

import (
	"github.com/leilei3167/cgin/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	time.Sleep(time.Second * 5)
	c.Json(200, "ok UserLoginController")
	return nil
}
