package rest

import (
	"context"
	"data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

type Input struct {
	Card         []string
	DiscountType []string
}
type Response struct {
	Message string      `json:"Message,omitempty"`
	Result  interface{} `json:"Result,omitempty"`
	Count   int         `json:"Count,omitempty"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func InitHackAPI() *gin.Engine {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.POST("/echo", func(c *gin.Context) {
		inputByte, _ := ioutil.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(inputByte))
	})

	router.GET("/googlemapjob", func(c *gin.Context) {
		// init google map client
		gClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyCOupeguSkCc8RDiLTkA94uk02Jzpq5xSo"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		ctx := context.Background()
		//get job
		jobs, err := data.GetJob()
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		// fmt.Println("jobs:", jobs)
		for _, job := range jobs {
			// ask google
			resp, err := gClient.Geocode(ctx, &maps.GeocodingRequest{Address: job.Address, Language: "zh-TW"})
			if err != nil {
				log.Println(err)
				continue
			}
			// fmt.Printf("%+v\n", resp)
			if resp == nil || len(resp) == 0 {
				log.Println("No data:", job.Id)
				continue
			}
			gData := resp[0]
			job.GoogleAddress = gData.FormattedAddress
			job.Latitude = gData.Geometry.Location.Lat
			job.Longutitude = gData.Geometry.Location.Lng
			// set
			data.Set(job)
		}

		c.JSON(http.StatusOK, Response{Count: len(jobs), Result: "OK"})
	})

	// input as /activity?input={"Card":["幸運PLUS鈦金卡"],"DiscountType":["食"]}
	router.GET("/activity", func(c *gin.Context) {
		inputStr, ok := c.GetQuery("input")
		if !ok {
			c.JSON(http.StatusBadRequest, Response{Message: "input"})
			return
		}
		input := Input{}
		err := json.Unmarshal([]byte(inputStr), &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			return
		}
		activities, err := data.GetActivity(input.Card, input.DiscountType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
			return
		}
		n := len(activities)
		weekday := strconv.Itoa(int(time.Now().Weekday()))
		fmt.Println("week", weekday)
		for i := 0; i < n; i++ {
			times := activities[i].Time
			activities[i].IsToday = strings.Contains(times, weekday)
		}
		c.JSON(http.StatusOK, activities)
	})

	return router
}
