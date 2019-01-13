package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	{
		group := r.Group("/biquge")
		group.GET("/", BiqugeIndex)
		group.GET("/:bookID", BiqugeRSS)
	}

	{
		group := r.Group("/dota")
		group.GET("/", DotaIndex)
		group.GET("/:listID", DotaRSS)
	}
	r.Run()
}
