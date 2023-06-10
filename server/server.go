package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
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

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies(nil)

	server.LoadHTMLGlob("../client/html/*.html")
	server.Static("/static", "../client/static")

	server.GET("/api/wiki", wiki_crawl)
	server.GET("/", func (c *gin.Context) {
		c.HTML(200, "home.html", "")
	})

	server.Run("0.0.0.0:4000")
}
