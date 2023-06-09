package main

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

func def(c *gin.Context) {
	c.IndentedJSON(200, "deeba")
}

func wiki_scrape(c *gin.Context) {
	url := c.Query("url")
	depth_ := c.Query("depth")

	depth, derr := strconv.ParseInt(depth_, 10, 64)
	if derr != nil {
		c.JSON(400, gin.H{"Error": "invalid depth value."})
		return
	}
	if depth < 1 || depth > 5 {
		c.JSON(400, gin.H{"Error": "invalid depth - expected an int between 1 and 5."})
		return
	}

	res, err := Scrape("https://en.wikipedia.org/wiki/" + url, int(depth))
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, res)
}

func main() {
	server := gin.Default()
	server.SetTrustedProxies(nil)

	server.LoadHTMLGlob("../client/html/*.html")
	server.Static("/static", "../client/static")

	server.GET("/api/wiki", wiki_scrape)
	server.GET("/", func (c *gin.Context) {
		c.HTML(200, "home.html", "")
	})

	server.Run("localhost:1337")
}
