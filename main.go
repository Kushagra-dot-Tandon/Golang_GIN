package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type request_json struct {
	APPID  int    `json:appid`
	Status string `json:status`
	User   string `json:user`
}

type awsconnect struct {
	Bucket_name string `json:"bucket_name"`
	Region_name string `json:"region_name"`
}

type AppProcess struct {
	gorm.Model
	AppID  int
	Status string
	User   string
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}

}

func connect_db() *gorm.DB {
	//Initalization of DATABASE
	db, err := gorm.Open("postgres", "user=postgres password=kush dbname=gorm sslmode=disable")
	CheckError(err)
	return db
}

func main() {

	//Initalization of GIN
	r := gin.Default()

	db := connect_db()
	defer db.Close()

	fmt.Printf("Database Connection Successful \n")

	r.POST("/update_process", func(c *gin.Context) {
		var data_json request_json
		c.BindJSON(&data_json)
		var data_to_database = &AppProcess{AppID: data_json.APPID, Status: data_json.Status, User: data_json.User}
		db.Create(data_to_database)
	})

	// READING DATA FROM JSONFILE
	r.POST("/readconfig", func(c *gin.Context) {
		file, err := os.Open("./config/config.json")
		CheckError(err)

		// Declaration for the json_data
		var aws awsconnect
		decoder := json.NewDecoder(file)

		err = decoder.Decode(&aws)
		CheckError(err)
		// fmt.Println(aws)
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
