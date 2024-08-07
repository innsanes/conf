package conf

import (
	"gopkg.in/yaml.v3"
)

type Yaml struct {
	YamlConf *YamlConf
	*kv[interface{}]
	conf *X
}

type YamlConf struct {
	FilePath string `conf:"filepath,default=config.yaml"`
}

func NewYaml(conf *X) *Yaml {
	return &Yaml{
		YamlConf: &YamlConf{},
		kv:       newKV[interface{}](),
		conf:     conf,
	}
}

func (y *Yaml) Parse() {
	// 判断文件是否存在, 如果存在则读取, 如果不存在就创建文件
	if !y.conf.fileExist(y.YamlConf.FilePath) {
		y.format()
		return
	}
	// 将文件中的yaml数据解析成map
	binaryData := y.conf.readFile(y.YamlConf.FilePath)
	var data map[string]interface{}
	err := yaml.Unmarshal(binaryData, &data)
	if err != nil {
		y.conf.log.Fatalf("yaml unmarshal err:%s", err)
	}
	y.yamlRecursiveParse(data, "")
}

func (y *Yaml) yamlRecursiveParse(data map[string]interface{}, prefix string) {
	for key, v := range data {
		// 判断value是否是map, 如果是继续递归, 如果不是, 存入KV中
		if subData, ok := v.(map[string]interface{}); ok {
			y.yamlRecursiveParse(subData, prefix+key+"_")
			continue
		}
		y.Set(prefix+key, v)
	}
}

func (y *Yaml) format() {
	// 将 conf 中的tree数据转成 map 并将其写入文件中
	// 1. 将tree数据转成map
	data := make(map[string]interface{})
	y.yamlRecursiveFormat(y.conf.argTree, data)
	// 2. 将map数据转成yaml
	binaryData, err := yaml.Marshal(data)
	if err != nil {
		y.conf.log.Fatalf("yaml marshal err:%s", err)
	}
	// 3. 将yaml写入文件
	y.conf.writeFile(y.YamlConf.FilePath, binaryData)
}

func (y *Yaml) yamlRecursiveFormat(tree *argTree, data map[string]interface{}) {
	for _, child := range tree.child {
		if len(child.child) == 0 {
			data[child.key] = child.value
			continue
		}
		subData := make(map[string]interface{})
		data[child.key] = subData
		y.yamlRecursiveFormat(child, subData)
	}
}
