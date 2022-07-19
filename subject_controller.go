package main

import "github.com/leilei3167/cgin/framework"

func SubjectListController(c *framework.Context) error {
	c.Json(200, "ok SubjectListController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	c.Json(200, "ok SubjectGetController")
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.Json(200, "ok SubjectUpdateController")
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.Json(200, "ok SubjectDelController")
	return nil
}
