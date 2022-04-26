package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type SinkModel struct {
	ID             uint           `gorm:"primary_key"`
	Name           string         `gorm:"column:name;UNIQUE"`
	ProfileID      string         `gorm:"column:profileid"`
	SinkType       string         `gorm:"column:type"`
	ConnectURL     string         `gorm:"column:connect_url"`
	MaxTasks       uint           `gorm:"column:max_tasks"`
	Config         postgres.Jsonb `gorm:"column:config"`
	State          string         `gorm:"column:state"`
	TrueState      string         `gorm:"column:true_state;default:'RUNNING'"`
	QuotaUsageStat postgres.Jsonb `gorm:"column:quota_usage_stats"`
}

type QuotaUsageStat struct {
	ProducerByteRate float64 `json:"producer_byte_rate"`
	ConsumerByteRate float64 `json:"consumer_byte_rate"`
}

func (SinkModel) TableName() string {
	return "kconnect_sinks"
}

//check_error
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("postgres", "host=snappyflow-stage-rds.cmebsmm3knzu.us-west-2.rds.amazonaws.com port=5432 user=archive dbname=archival password=archive123 sslmode=disable")
	CheckError(err)
	return db
}

func queryPrometheus(profileRecord SinkModel, query string) float64 {
	var totalQuotaSum float64 = 0.0

	client, err := api.NewClient(api.Config{
		Address: "http://54.186.248.93:32031",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	for _, val := range result.(model.Vector) {
		totalQuotaSum += float64(val.Value)
	}
	return totalQuotaSum
}

func main() {
	// Connect to Database
	Db := initDatabase()
	//Close the Database after main is over
	defer Db.Close()

	var profileRecords []SinkModel
	var quotaStat QuotaUsageStat

	Db.Model(&SinkModel{}).Find(&profileRecords)
	for _, profileRecord := range profileRecords {
		query1 := fmt.Sprintf("sum by (instance) (avg_over_time(kafka_producer_client_topic_byte_rate{app=\"cp-kafka-rest\",topic=~\"[a-zA-Z]+-%s\"}[30m]))", profileRecord.ProfileID)
		totalQuotaSum := queryPrometheus(profileRecord, query1)
		quotaStat.ProducerByteRate = totalQuotaSum
		jsonByte, err := json.Marshal(&quotaStat)
		if err != nil {
			fmt.Printf("Error while marshal quotaStat for Profile %s", profileRecord.ProfileID)
		}
		if err = Db.Model(&SinkModel{}).Where("profileid = ?", profileRecord.ProfileID).Update("quota_usage_stats", postgres.Jsonb{RawMessage: json.RawMessage(jsonByte)}).Error; err != nil {
			CheckError(err)
		}
		fmt.Printf("Done ")
		time.Sleep(2 * time.Second)

		query2 := fmt.Sprintf("avg_over_time(kafka_consumer_client_topic_bytes_consumed_rate{app=\"es-kafka-connect\",topic=~\"[a-zA-Z]+-%s\"}[30m])", profileRecord.ProfileID)
		totalQuotaSum = queryPrometheus(profileRecord, query2)
		quotaStat.ConsumerByteRate = totalQuotaSum
		jsonByte2, err := json.Marshal(&quotaStat)
		if err != nil {
			fmt.Printf("Error while marshal quotaStat for Profile %s", profileRecord.ProfileID)
		}
		if err = Db.Model(&SinkModel{}).Where("profileid = ?", profileRecord.ProfileID).Update("quota_usage_stats", postgres.Jsonb{RawMessage: json.RawMessage(jsonByte2)}).Error; err != nil {
			CheckError(err)
		}
		time.Sleep(2 * time.Second)

		if err := json.Unmarshal(profileRecord.QuotaUsageStat.RawMessage, &quotaStat); err != nil {
			CheckError(err)
		}
		fmt.Println("", quotaStat.ConsumerByteRate)
	}

}
