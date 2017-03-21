package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"io"

	"../model"
	"github.com/gin-gonic/gin"
	"github.com/johnhjwsosd/redis-operation/redisoper"
)

//GetRouter ...
func GetRouter() *gin.Engine {
	router := gin.New()
	router.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to apigateway")
	})
	router.POST("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome to apigateway")
	})

	router.GET("/api/:ms", func(c *gin.Context) {
		msName := c.Param("ms")
		msg := fmt.Sprint("please input microserver Api----Get ", msName)
		c.String(http.StatusOK, msg)
	})
	router.POST("/api/:ms", func(c *gin.Context) {
		msName := c.Param("ms")
		msg := fmt.Sprint("please input microserver Api----Post ", msName)
		c.String(http.StatusOK, msg)
	})
	router.GET("/usercenter/login", loginMethods)
	router.GET("/usercenter/logout", logoutMethods)
	router.POST("/usercenter/login", loginMethods)
	router.POST("/usercenter/logout", logoutMethods)

	router.GET("/api/:ms/:api", requestHandles)
	router.POST("/api/:ms/:api", requestHandles)
	return router
}
func loginMethods(c *gin.Context) {
	urlHost := "http://192.168.1.178:9999/login?" + c.Request.URL.RawQuery
	data := c.Request.Body
	resModel := requestMethods(urlHost, data, "POST")
	fmt.Println(resModel)
	switch resModel.StatusCode {
	case 0:
		c.JSON(500, resModel)
	case 4:
	case 5:
	case 1000:
		c.JSON(200, resModel.Info)
	}
	return
}

func logoutMethods(c *gin.Context) {
	urlHost := "http://192.168.1.189:8080/user/logout?" + c.Request.URL.RawQuery
	data := c.Request.Body
	resModel := requestMethods(urlHost, data, "POST")

	switch resModel.StatusCode {
	case 0:
		c.JSON(500, resModel)
	case 4:
	case 5:
	case 1000:
		c.JSON(200, resModel.Info)
	}
	return
}

func requestHandles(c *gin.Context) {
	params := getParams(c)

	temp := getInfo(params["MS"], params["API"])

	if temp.MSHost == "" {
		c.JSON(404, "server is not found")
		return
	}
	if temp.APIMethods == "" {
		c.JSON(404, "API is not found")
		return
	}

	urlHost := temp.MSHost + "/" + temp.MsName + "/" + temp.APIName + "?" + c.Request.URL.RawQuery
	fmt.Println(urlHost)
	data := c.Request.Body

	if temp.TokenMethods == "1" {
		//请求AUTH
		resModel := requestMethods(urlHost, data, "POST")
		if resModel.StatusCode == 1000 {
			c.JSON(200, resModel)
			return
		}
		c.JSON(500, "server occur fatal")
		return

	}
	//直接请求微服务
	resModel := requestMethods(urlHost, data, "POST")
	if resModel.StatusCode == 1000 {
		c.JSON(200, resModel)
		return
	}
	c.JSON(500, "server occur fatal")
	return
}

func getParams(c *gin.Context) map[string]string {
	msName := c.Param("ms")
	apiName := c.Param("api")
	r := c.Request
	r.ParseForm()
	params := make(map[string]string)
	params["MS"] = msName
	params["API"] = apiName
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
	return params
}

func getRedisData(key string) (interface{}, error) {
	hostString := "192.168.1.91:6379"
	authString := "123"
	redis := redisoper.NewRedis(hostString, authString)
	res, err := redis.Smembers(key)
	return res, err
}

func getMS(msName string) *model.MsInfo {
	res, err := getRedisData("MS:" + msName)
	handleError(err)
	ms := &model.MsInfo{}
	resList := res.([]string)
	if len(resList) != 0 {
		err = json.Unmarshal([]byte(resList[0]), ms)
	}
	return ms
}

func getInfo(msName, apiName string) *model.RequestModel {
	msModel := getMS(msName)
	res, err := getRedisData("API:" + msName)
	handleError(err)
	resList := res.([]string)
	result := &model.RequestModel{MsName: msName, MSHost: msModel.Host, APIName: apiName}
	for _, v := range resList {
		api := &model.APIInfo{}
		err := json.Unmarshal([]byte(v), api)
		handleError(err)
		if api.APIName == apiName {
			result.APIMethods = api.APIMethods
			result.TokenMethods = api.TokenMethods
		}
	}
	return result
}

func requestMethods(url string, data io.Reader, requestType string) *model.ComResult {
	var res []byte
	if requestType == "POST" {
		res = postRequest(url, "application/json", data)
	} else {
		res = getRequest(url)
	}
	resModel := &model.ComResult{}
	err := json.Unmarshal(res, resModel)
	handleError(err)
	return resModel
}

func handleError(err error) {
	if err != nil {
		fmt.Println("fatal occur :", err)
	}
}

func getRequest(url string) []byte {
	res, err := http.Get(url)
	handleError(err)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body
}

func postRequest(url, bodyType string, data io.Reader) []byte {
	res, err := http.Post(url, bodyType, data)
	if res == nil {
		return nil
	}
	defer res.Body.Close()
	handleError(err)
	body, _ := ioutil.ReadAll(res.Body)
	return body
}
