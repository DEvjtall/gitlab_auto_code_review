package interactive

import (
	"ark/muscle/config"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func gitlabComment(projectID string, comment string, commitID string) error {
	//初始化日志
	log := setupLog()
	log_file := log.Out.(*os.File)
	defer log_file.Close()
	// ---------------------------
	gitlab_url, _ := config.ReadConfig("gitlab", "GITLAB_URL")
	private_token, _ := config.ReadConfig("gitlab", "PRIVATE_TOKEN")
	// 创建一个缓冲区来保存表单数据
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// 设置 comment note
	err := writer.WriteField("note", comment)
	if err != nil {
		log.Error("设置comment note失败")
		return err
	}
	writer.Close()
	url := fmt.Sprintf(gitlab_url+"%s/repository/commits/%s/comments", projectID, commitID)
	fmt.Println(url)
	// 创建 comment POST 请求体
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		log.Error("comment API 创建请求失败:", err)
		return err
	}
	// 设置请求头
	req.Header.Set("PRIVATE-TOKEN", private_token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送 comment 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("comment API 发送请求失败")
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("comment API 读取响应失败")
		return err
	}
	log.Infof("响应状态码：%d", resp.StatusCode)
	log.Infof("响应体：%s", string(body))
	return nil
}
