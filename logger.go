package conf

import "fmt"

func resultHandler(result *ParseResult) {
	if result.Err != nil {
		panic(fmt.Sprintf("config parse fail: %s", result.Err))
	}
	for _, config := range result.configs {
		fmt.Println(fmt.Sprintf("-%s:%v, default:%s, usage:%s", config.Key, config.Value, config.Default, config.Usage))
	}
}
