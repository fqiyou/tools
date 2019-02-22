
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fqiyou/tools/foo/system"
	"github.com/fqiyou/tools/foo/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var (
	// 命令行参数
	listenAddr  = flag.String("web.listen-port", "9102", "An port to listen on for web interface and telemetry.")
	action_shell = flag.String("action_shell", "/Users/chaoyang/tmp/a.sh", "触发执行shell脚本名称")
	gitlab_project_id = flag.String("gitlab_project_id", "626", "gitlab项目id")
	logger = util.Log
	last_trigger_info map[string]interface{}
)

func execScript(shell_file string) {
	status,err := system.ExecShellScript(shell_file)
	if err != nil {
		logger.Error("执行脚本报错")
		logger.Error(err)
		logger.Error(status)
	}
}


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		//c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}




func Trigger(c *gin.Context) {

	//logger.Info("c.Request.Method: %v", c.Request.Method)
	//logger.Info("c.Request.ContentType: %v", c.ContentType())
	//logger.Info("c.Request.Body: %v", c.Request.Body)
	//c.Request.ParseForm()
	//logger.Info("c.Request.Form: %v", c.Request.PostForm)
	//for k, v := range c.Request.PostForm {
	//	logger.Info("k:%v\n", k)
	//	logger.Info("v:%v\n", v)
	//}
	//logger.Info("c.Request.ContentLength: %v", c.Request.ContentLength)
	data, _ := ioutil.ReadAll(c.Request.Body)
	logger.Info("c.Request.GetBody: ", string(data))

	var trigger_info map[string]interface{}


	err := json.Unmarshal([]byte(string(data)), &trigger_info)
	last_trigger_info = trigger_info

	if err == nil {

		project_id := fmt.Sprintf("%d", int(trigger_info["project_id"].(float64)))
		if project_id == *gitlab_project_id {
			execScript(*action_shell)
		}
	} else {
		logger.Error(err)

	}
	c.JSON(200, map[string]interface{}{"status": 200, "status_info": "success", "data": ""})



}


func main() {
	flag.Parse()
	gin.SetMode(gin.DebugMode)


	router :=gin.Default()
	router.Use(CORSMiddleware())
	router.Any("/", func(context *gin.Context) {
		context.String(http.StatusOK,"hello")
	})
	v1 := router.Group("/v1")
	{
		v1.Any("/", func(context *gin.Context) {
			context.String(http.StatusOK,"v0.0.1")
		})
		v1.POST("/trigger",Trigger)
		v1.Any("/trigger_info", func(context *gin.Context) {
			context.JSON(200, last_trigger_info)
		})



	}
	router.Run(":" + *listenAddr)

	util.Log.Info("Starting Server at http://localhost:", *listenAddr)
}

