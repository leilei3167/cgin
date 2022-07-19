package main

import "github.com/leilei3167/cgin/framework"

func UserLoginController(c *framework.Context) error {
	c.Json(200, "ok UserLoginController")
	return nil
}
