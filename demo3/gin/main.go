package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/user/:name", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": c.Param("name"),
			"msg":  "success",
		})
	})

	r.Run(":8080")
}
