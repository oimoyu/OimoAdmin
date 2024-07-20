package config

import (
	"encoding/json"
	"fmt"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"os"
)

func LoadConfig() {
	path := functions.GetEnvDir() + "/config.json"
	config := new(_type.ConfigStruct)

	// 读取JSON文件
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// 反序列化JSON到config
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	g.Config = config

	fmt.Printf("Config loaded: %+v\n", config)
}
