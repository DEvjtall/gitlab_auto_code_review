package interactive

import (
	"ark/muscle/arkai"
	"ark/muscle/config"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"golang.org/x/net/context"
	"os"
	"strings"
)

type ArkClientConfig struct {
	APIKEY  string
	BaseURL string
	Region  string
}

// 初始化日志
func setupLog() *logrus.Logger {
	// 日志初始化
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // 完整的时间戳，可以看到日志发生的准确时间
	})
	file, _ := os.OpenFile("ark.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	log.SetOutput(file) // 在这里还不能够把日志输出关闭，在这里关闭的话程序的报错日志就写入不到文件里面的了
	return log
}

// 获取上传文件中最新的commit ID以便写到 gitlab 评论里面
func getID(filePath string, prefix string) (string, error) {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix)), nil
		}
	}
	return "", nil
}

// 读取 git_diff 内容
func readGitDiff(filePath string) (string, error) {
	diffent_content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(diffent_content), nil
}

// 设置上传文件的 handler
func UploadHandler(c *gin.Context) {
	// log 初始化
	log := setupLog()
	log_file := log.Out.(*os.File)
	defer log_file.Close()
	// -----------------------------

	// 尝试打开上传文件
	diff_file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "无法获得上传的文件"})
		log.Errorf("无法获得上传的文件：%v", err)
		return
	}
	// 将上传的diff文件保存到本地然后读取
	upload_file := "upload_file" + diff_file.Filename
	file_path := "UpLoadDir" + "/" + upload_file
	if err := c.SaveUploadedFile(diff_file, file_path); err != nil {
		c.JSON(500, gin.H{"error": "无法保存上传的文件"})
		log.Errorf("无法保存上传的文件：%v", err)
		return
	}

	// 直接读取文件内容（这个后面发现是读取不了的，要把上传的文件保存到本地之后读取）
	content, err := readGitDiff(file_path)
	if err != nil {
		c.JSON(400, gin.H{"error": "无法读取上传的文件"})
		log.Errorf("无法读取上传的文件：%v", err)
		return
	}
	// 调用 getID()
	commit_id, err := getID(file_path, "commitID:")
	project_id, err := getID(file_path, "projectID:")
	// 向 ark_init 发送 git_diff 内容，然后调用 AI 接口 code review
	if err := ark_init(project_id, commit_id, content); err != nil {
		c.JSON(500, gin.H{"error": "无法初始化Ark"})
		log.Errorf("无法初始化Ark：%v", err)
		return
	}
}

// ark AI 审核阶段
func NewArkClient(config ArkClientConfig) (*arkruntime.Client, error) {
	if config.APIKEY == "" {
		return nil, fmt.Errorf("API Key 不能为空")
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("基础 URL 不能为空")
	}
	if config.Region == "" {
		return nil, fmt.Errorf("区域不能为空")
	}
	client := arkruntime.NewClientWithApiKey(
		config.APIKEY,
		arkruntime.WithRegion(config.BaseURL),
		arkruntime.WithRegion(config.Region),
	)
	return client, nil
}

func ark_init(projectID string, commitID string, content string) error {
	// 创建ark日志初始化
	log := setupLog()
	log_file := log.Out.(*os.File)
	defer log_file.Close()
	// ------------------------
	// 给 ark_config 读取配置
	ark_key, _ := config.ReadConfig("default", "API_KEY")
	arkconfig := ArkClientConfig{
		APIKEY:  ark_key,
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Region:  "cn-guangzhou",
	}
	// 构建请求体
	client, err := NewArkClient(arkconfig)
	if err != nil {
		log.Error("创建客户端初始化失败")
		return err
	}
	ctx := context.Background()
	comments, _ := arkai.SendToArk(ctx, client, content)
	log.Info("Ark 初始化成功")

	// 发送gitlab评论
	err = gitlabComment(projectID, *comments, commitID)
	if err != nil {
		log.Error("gitlab发送代码审查评论失败")
	}

	return nil
}
