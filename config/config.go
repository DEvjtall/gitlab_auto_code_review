package config

import (
	"github.com/Unknwon/goconfig"
	"github.com/sirupsen/logrus"
	"os"
)

func ReadConfig(section string, key string) (string, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // 完整的时间戳，可以看到日志发生的准确时间
	})
	file, err := os.OpenFile("ark.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Errorf("打开配置相关错误日志出现问题：%v", err)
		return "", err
	}
	log.SetOutput(file)

	cfg, err := goconfig.LoadConfigFile("config/config.ini")
	if err != nil {
		log.Errorf("打开配置文件出现问题，%v", err)
		return "", err
	}

	apiKey, err := cfg.GetValue(section, key)
	if err != nil {
		log.Errorf("获取配置文件中的KEY值失败：%v", err)
		return "", err
	}

	return apiKey, nil
}
