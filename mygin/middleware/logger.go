package middleware

import (
	"log"
	mygin "myWeb/mygin"
	"time"
)

func Logger() mygin.HandlerFunc {
	return func(c *mygin.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
