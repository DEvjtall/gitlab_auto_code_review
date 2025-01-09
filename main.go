package main

import (
	"ark/muscle/interactive"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func setuplog() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // 完整的时间戳，可以看到日志发生的准确时间
	})
	file, _ := os.OpenFile("ark.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	log.SetOutput(file)
	return log
}

func main() {
	//// 日志初始化
	log := setuplog()
	log_file := log.Out.(*os.File)
	defer log_file.Close()
	// -----------------------------
	r := gin.Default()
	// 处理文件上传的端口
	r.POST("/upload", interactive.UploadHandler)
	// 自定义端口
	port := 7777
	log.Infof("服务器启动，正在监听：%d...", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}

}
