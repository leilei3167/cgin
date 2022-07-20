# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/leilei3167/cgin/framework/gin"
	"github.com/leilei3167/cgin/framework/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
