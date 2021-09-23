package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Define the type of the request_JSON
type request_json struct {
	APPID  int    `json:appid`
	Status string `json:status`
	User   string `json:user`
}

//Define the type of the AWS
type awsconnect struct {
	Bucket_name string `json:"bucket_name"`
	Region_name string `json:"region_name"`
}

//Define the type of the database
type AppProcess struct {
	gorm.Model
	AppID  int
	Status string
	User   string
}

//check_error
func CheckError(err error) {
	if err != nil {
		panic(err)
	}

}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("postgres", "user=postgres password=kush dbname=gorm sslmode=disable")
	CheckError(err)
	return db
}

func main() {

	//Initalization of GIN
	r := gin.Default()
	// Connect to Database
	db := initDatabase()
	//Close the Database after main is over
	defer db.Close()

	//Update Process => To Update data onto Database
	r.POST("/update_process", func(c *gin.Context) {
		var data_json request_json
		c.BindJSON(&data_json)
		var data_to_database = &AppProcess{AppID: data_json.APPID, Status: data_json.Status, User: data_json.User}
		db.Create(data_to_database)
	})

	//Query Time => To check the database onto the query
	r.GET("time_query/:hour", func(c *gin.Context) {
		// time_select := c.Param("hour")
		//Intialization the data wrt to our database
		var data []AppProcess
		// db.find to find and query all the content .......
		// db.Find(&data)
		//let us suppose the expiry time is one day before: yesterday so => dt.AddDate(0,0,-1)
		db.Where("created_at < ?", time.Now().AddDate(0, 0, -1)).Find(&data)
		//  iterative onto the database and get all the fields or jobs having timespan less than 2 hours
		for _, u := range data {
			fmt.Println(u.AppID)
			// database_time := u.CreatedAt
			// diff := dt.Sub(database_time)
			// fmt.Println(int(diff.Hours()))
			// if int(diff.Hours()) <= 2 {
			// 	fmt.Println("AppId", u.AppID, "Process", u.Status)
			// }
		}
	})

	// READING DATA FROM JSONFILE
	r.POST("/readconfig", func(c *gin.Context) {
		file, err := os.Open("./config/aws.json")
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
			id := c.Param("id")
			c.JSON(200, gin.H{
				"route":        true,
				"resources_id": id,
			})
		})

		id := resources.Group("/:id")
		{
			subresource := id.Group("/subresource")
			{
				subresource.GET("/:newid", func(c *gin.Context) {
					new_id := c.Param("newid")
					c.JSON(200, gin.H{
						"route":       true,
						"subresource": "Working Perfectly_In_Sub_Route",
						"data_new_id": new_id,
					})
				})
			}

		}
	}
	r.Run()
}
