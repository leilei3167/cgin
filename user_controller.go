package main

import (
	"context"
	"errors"
	"github.com/leilei3167/cgin/framework"
	"net/http"
	"time"
)

func UserLoginController(c *framework.Context) error {
	ctx, cancel := context.WithTimeout(c.BaseContext(), time.Second*3)
	defer cancel()

	//执行处理逻辑
	if err := login(ctx); err != nil {
		c.Json(http.StatusBadRequest, err.Error())
		return nil
	}

	c.Json(200, "success!")

	return nil
}

func login(ctx context.Context) error {
	done := make(chan bool)

	go func() {
		//处理
		time.Sleep(time.Second * 10)
		done <- true
	}()

	select {
	case <-ctx.Done():
		return errors.New("time out")

	case <-done:
		return nil
	}

}
