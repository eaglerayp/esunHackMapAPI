package rest

import (
	"data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Input struct {
	Card         []string
	DiscountType []string
}
type Response struct {
	Message string      `json:"Message,omitempty"`
	Result  interface{} `json:"Result,omitempty"`
}

func InitHackAPI() *gin.Engine {
	router := gin.Default()

	router.POST("/echo", func(c *gin.Context) {
		inputByte, _ := ioutil.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(inputByte))
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
