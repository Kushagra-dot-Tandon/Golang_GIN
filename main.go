package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Kushagra_Maple_Labs")
	r := gin.Default()

	r.POST("/awsconnect", func(c *gin.Context) {
		data := c.Request.Body
		value, err := ioutil.ReadAll(data)
		if err != nil {
			fmt.Println(err.Error())
		}
		c.JSON(200, gin.H{
			"message": "Hello Kushagra_Maple_Labs AWS Connect",
			"data":    string(value),
		})
	})
	// STACKOVERFLOW PROBLEM SOLUTION -> SUBROUTING In GOLANG

	resources := r.Group("/resources")
	{
		resources.GET("/:id", func(c *gin.Context) {
			res_id := c.Param("id")
			c.JSON(200, gin.H{
				"route":        true,
				"resources_id": res_id,
			})
		})

		id := resources.Group("/:id")
		{
			subresource := id.Group("/subresource")
			{
				subresource.GET("/:newid", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"route":       true,
						"subresource": "Working Perfectly_In_SubRoute",
					})
				})
			}

		}
	}
	r.Run()
}
