package router

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johnhjwsosd/redis-operation/redisoper"
)

//GetRouter ...
func GetRouter() *gin.Engine {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to apigateway")
	})
	router.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to apigateway")
	})

	router.GET("/:ms", func(c *gin.Context) {
		msName := c.Param("ms")
		msg := fmt.Sprint("please input microserver Api----Get ", msName)
		c.String(http.StatusOK, msg)
	})
	router.POST("/:ms", func(c *gin.Context) {
		msName := c.Param("ms")
		msg := fmt.Sprint("please input microserver Api----Post ", msName)
		c.String(http.StatusOK, msg)
	})

	router.GET("/:ms/:api", requestHandles)
	router.POST("/:ms/:api", requestHandles)
	return router
}

func requestHandles(c *gin.Context) {
	msName := c.Param("ms")
	apiName := c.Param("api")
	r := c.Request
	r.ParseForm()
	params := make(map[string]interface{})
	params["welcome"] = msName
	params["msg"] = apiName
	if r.Method == "GET" {
		for k, v := range r.Form {
			if c.Query(k) != "" {
				params[k] = v[0]
			}
		}
	} else {
		for k, v := range r.PostForm {
			params[k] = v[0]
		}
	}
	res, err := getRedisData(msName, apiName)
	if err != nil {
		c.JSON(500, err)
		return
	}
	params["data"] = res
	c.JSON(200, params)
	return
}

func getRedisData(msName, apiName string) (interface{}, error) {
	hostString := "192.168.1.41:6379"
	authString := "123"
	redis := redisoper.NewRedis(hostString, authString)
	pool := redis.NewPool()
	res, err := redis.GetData(pool, msName, "set")
	return res, err
}
