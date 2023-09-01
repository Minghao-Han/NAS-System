package Utils

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

func YmlReader(filePath string, configPath string) interface{} {
	data := &map[string]interface{}{}
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read the YAML file: %v", err)
	}

	// 创建一个空的 map[string]interface{} 对象，用于解析 YAML 数据
	// 解析 YAML 数据到 map 中
	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	path := strings.Split(configPath, ":")
	value := (*data)[path[0]]
	for i := 1; i < len(path); i++ {
		if value != nil {
			key := path[i]
			value = value.(map[string]interface{})[key]
		}
	}
	return value
}
