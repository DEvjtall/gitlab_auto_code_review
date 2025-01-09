package arkai

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"os"
)

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

func SendToArk(ctx context.Context, client *arkruntime.Client, codeContent string) (*string, error) {
	// 初始化日志文件
	log := setupLog()
	file_log := log.Out.(*os.File)
	defer file_log.Close()
	// -----------------------------
	req := model.BotChatCompletionRequest{
		BotId: "bot-20241215214415-ct7lv", // 智能体的ID
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("从现在开始，你是公司的代码审查员，有着极其非常丰富的编程经验，熟悉各类开发语言、框架的，根据以下标准对变更的代码审查标准如下： 1.Bug 检测：查找潜在 Bug 或逻辑错误, 并定义风险级别高、中、低。 2.代码质量：查找不必要的复杂性或冗余代码。 3.安全实践：检查漏洞和硬编码密钥。 4.性能优化：提出更改建议以提高效率，指出问题代码。 5.语法和样式：查找语法错误和与约定的偏差。 代码审查完给出本次审查报告，审查报告需要使用markdown格式输出， 注意输出内容的美观和可阅读性， 报告格式如下： 一、总评： 首先你总评下本次代码变更的内容，如果有发现bug需要特别说明发现有多少个bug, 要求表达简洁不超过120个字 二、审查详情：需要根据上述6条代码审查标准分别点评，每个点评如有发现问题代码，需要贴出具体的问题代码，并提供清晰、可操作的优化方案 。 如果以上某个指标没问题， 可以简洁回复，比如：语法和样式：暂时未发现明显优化点, 每条审查标准的点评需要使用空行分开便于阅读。 以下是代码变更:"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(codeContent),
				},
			},
		},
	}
	resp, err := client.CreateBotChatCompletion(ctx, req)
	if err != nil {
		log.Errorf("发送请求失败：%v", err)
		return nil, fmt.Errorf("发送请求失败：%v", err)
	}
	log.Infof("请求成功：%v", resp)
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content.StringValue != nil {
		fmt.Println(*resp.Choices[0].Message.Content.StringValue)
		return resp.Choices[0].Message.Content.StringValue, nil
	}
	return nil, fmt.Errorf("没有收到生成的文本")
}
