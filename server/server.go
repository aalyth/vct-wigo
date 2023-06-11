package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/JGLTechnologies/gin-rate-limit"
)

func def(c *gin.Context) {
	c.IndentedJSON(200, "deeba")
}

func wiki_crawl(c *gin.Context) {
	url := c.Query("url")
	depth_ := c.Query("depth")

	depth, derr := strconv.ParseInt(depth_, 10, 64)
	if derr != nil {
		c.JSON(400, gin.H{"Error": "invalid depth value."})
		return
	}
	if depth < 1 || depth > 3 {
		c.JSON(400, gin.H{"Error": "invalid depth - expected an int between 1 and 5."})
		return
	}

	res := Scrape(url, int(depth))
	c.JSON(200, res)
}

func healthcheck(c *gin.Context) {
	c.Status(200)
}


func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies(nil)

	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  2 * time.Second,
		Limit: 1,
	})
	limiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc: keyFunc,
	})

	server.LoadHTMLGlob("../client/html/*.html")
	server.Static("/static", "../client/static")

	server.GET("/api/wiki", limiter, wiki_crawl)
	server.GET("/hc", healthcheck)
	server.GET("/", func (c *gin.Context) {
		c.HTML(200, "home.html", "")
	})

	server.Run(":4000")
}
