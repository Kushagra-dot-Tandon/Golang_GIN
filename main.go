package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	r.POST("/readconfig", func(c *gin.Context) {
		file, err := os.Open("./config/config.json")
		if err != nil {
			fmt.Print(err)
		}

		type awsconnect struct {
			Bucket_name string `json:"bucket_name"`
			Region_name string `json:"region_name"`
		}

		var aws awsconnect
		decoder := json.NewDecoder(file)

		err = decoder.Decode(&aws)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(aws)
		c.JSON(200, gin.H{
			"message":     "Hello Kushagra_Maple_Labs AWS Connect",
			"bucket_name": aws.Bucket_name,
			"region":      aws.Region_name,
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
