package Utils

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	data map[string]interface{}
}

var configuration *Configuration

func DefaultConfigReader() *Configuration {
	if configuration == nil {
		configuration = &Configuration{}
		configuration.initial()
	}
	return configuration
}

func (config *Configuration) initial() {
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("Failed to read the YAML file: %v", err)
	}

	// 创建一个空的 map[string]interface{} 对象，用于解析 YAML 数据
	// 解析 YAML 数据到 map 中
	err = yaml.Unmarshal(yamlFile, &config.data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	// 访问部分数据
	//serverHost := data["server"].(map[interface{}]interface{})["host"]
}

func (config *Configuration) Get(configPath string) interface{} { //path格式：database:password
	path := strings.Split(configPath, ":")
	value := config.data[path[0]]
	for i := 1; i < len(path); i++ {
		if value != nil {
			key := path[i]
			value = value.(map[string]interface{})[key]
		}
	}
	return value
}
