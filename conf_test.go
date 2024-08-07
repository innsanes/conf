package conf_test

import (
	"conf"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

// test 函数会有默认的flag传入参数和flag.Parse()
// 所以在测试中需要在flag.Parse()之后再设置一下 os.Args, 也就是每个Test开始之前

func init() {
	// 在测试开始前删除test文件夹下面的所有  .yaml 文件
	os.RemoveAll("test")
	// 创建 test 文件夹存放测试yaml
	os.Mkdir("test", os.ModePerm)
}

type TestFlagTypeStruct struct {
	String        string `conf:"string"`                     // string
	Int           int    `conf:"int"`                        // int
	Bool          bool   `conf:"bool"`                       // bool
	StringDefault string `conf:"string_default,default=123"` // string default
	IntDefault    int    `conf:"int_default,default=123"`    // int default
	BoolDefault   bool   `conf:"bool_default,default=true"`  // bool default
}

// 测试 string, int, bool 三个常见类型
func TestFlagType(t *testing.T) {
	os.Args = []string{"", "-t_string=1", "-t_int=2", "-t_bool=true"}
	var x = conf.New()
	flag := conf.NewFlag(x)
	s := &TestFlagTypeStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterSource(flag)
	x.Parse()
	x.PrintResult()
	assert.Equal(t, "1", s.String)
	assert.Equal(t, 2, s.Int)
	assert.Equal(t, true, s.Bool)
	assert.Equal(t, "123", s.StringDefault)
	assert.Equal(t, 123, s.IntDefault)
	assert.Equal(t, true, s.BoolDefault)
}

type TestFlagNestedStruct struct {
	TestFlagNestedChild                     // 匿名继承
	Child               TestFlagNestedChild // 默认
	ChildBBBBB          TestFlagNestedChild `conf:"c"` // 修改名称
	ChildIgnore         TestFlagNestedChild `conf:"-"` // 忽略
}

type TestFlagNestedChild struct {
	Argument string `conf:"a"`
}

// 嵌套结构体, 匿名继承, 非匿名继承, 忽略, 以及更改结构体名称
func TestFlagNested(t *testing.T) {
	os.Args = []string{"", "-t_a=1", "-t_child_a=2", "-t_c_a=3"}
	var x = conf.New()
	flag := conf.NewFlag(x)
	s := &TestFlagNestedStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterSource(flag)
	x.Parse()
	x.PrintResult()
	assert.Equal(t, "1", s.Argument)
	assert.Equal(t, "2", s.Child.Argument)
	assert.Equal(t, "3", s.ChildBBBBB.Argument)
	assert.Equal(t, "", s.ChildIgnore.Argument)
}

type TestYamlGenStruct struct {
	TestYamlStructNest TestYamlStructNest `conf:"struct"`
}

type TestYamlStructNest struct {
	Name  string `conf:"name,default=nest"`
	Value int    `conf:"value,default=1024"`
}

// yaml 文件生成, 包括默认值
func TestYamlGen(t *testing.T) {
	var filepath = "test/test_gen.yaml"
	os.Args = []string{"", "-t_name=new", "-t_value=19223", "-yaml_filepath=" + filepath}
	var x = conf.New()
	flag := conf.NewFlag(x)
	y := conf.NewYaml(x)
	s := &TestYamlGenStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterConfWithName("yaml", y.YamlConf)
	x.RegisterSource(flag)
	x.RegisterSource(y)
	x.Parse()
	x.PrintResult()
	// 打开文件 验证文件内容是否与预期的一致
	buf, err := os.ReadFile(filepath)
	assert.Nil(t, err)
	ns := map[string]interface{}{}
	err = yaml.Unmarshal(buf, &ns)
	assert.Nil(t, err)
	assert.Equal(t, "nest", ns["t"].(map[string]interface{})["struct"].(map[string]interface{})["name"])
	assert.Equal(t, "1024", ns["t"].(map[string]interface{})["struct"].(map[string]interface{})["value"])
}

type TestYamlStruct struct {
	// yaml 标签是为了生成 yaml 文件, 测试 conf 的读取功能
	// 实际应用中不需要 yaml 标签
	TestYamlNest TestYamlNest `conf:"struct" yaml:"struct"`
}

type TestYamlNest struct {
	Name  string `conf:"name,default=nest" yaml:"name"`
	Value int    `conf:"value,default=1024" yaml:"value"`
}

// 测试 使用yaml 设置参数值, conf能否正确识别
func TestYaml(t *testing.T) {
	var filepath = "test/test.yaml"
	os.Args = []string{"", "-yaml_filepath=" + filepath}
	var x = conf.New()
	flag := conf.NewFlag(x)
	y := conf.NewYaml(x)
	s := &TestYamlStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterConfWithName("yaml", y.YamlConf)
	x.RegisterSource(flag)
	x.RegisterSource(y)

	data := map[string]interface{}{}
	expect := &TestYamlStruct{
		TestYamlNest: TestYamlNest{
			Name:  "custom",
			Value: 11011,
		},
	}
	data["t"] = expect
	marshal, err := yaml.Marshal(&data)
	assert.Nil(t, err)
	// 将yaml写入文件
	err = os.WriteFile(filepath, marshal, os.ModePerm)
	assert.Nil(t, err)

	x.Parse()
	x.PrintResult()
	assert.Equal(t, expect.TestYamlNest.Value, s.TestYamlNest.Value)
	assert.Equal(t, expect.TestYamlNest.Name, s.TestYamlNest.Name)
}

type TestFlagYamlStruct struct {
	TestFlagYamlStructNest TestFlagYamlStructNest `conf:"struct" yaml:"struct"`
}

type TestFlagYamlStructNest struct {
	Name  string `conf:"name,default=nest" yaml:"name"`
	Value int    `conf:"value,default=1024" yaml:"value"`
}

// 测试 yaml 和 flag 的优先级
func TestFlagYaml(t *testing.T) {
	var filepath = "test/test_priority.yaml"
	os.Args = []string{"", "-t_struct_name=flag-name", "-t_struct_value=19000", "-yaml_filepath=" + filepath}
	var x = conf.New()
	flag := conf.NewFlag(x)
	y := conf.NewYaml(x)
	s := &TestFlagYamlStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterConfWithName("yaml", y.YamlConf)
	x.RegisterSource(flag)
	x.RegisterSource(y)

	data := map[string]interface{}{}
	expect := &TestFlagYamlStruct{
		TestFlagYamlStructNest: TestFlagYamlStructNest{
			Name:  "yaml-name",
			Value: 1111,
		},
	}
	data["t"] = expect
	marshal, err := yaml.Marshal(&data)
	assert.Nil(t, err)
	// 将yaml写入文件
	err = os.WriteFile(filepath, marshal, os.ModePerm)
	assert.Nil(t, err)

	x.Parse()
	x.PrintResult()
	assert.Equal(t, "flag-name", s.TestFlagYamlStructNest.Name)
	assert.Equal(t, 19000, s.TestFlagYamlStructNest.Value)
}

type TestFuncStruct struct {
	Name string `conf:"name"`
}

// 非 flag 或者 yaml 这类官方配置参数的方式
// 通过函数调用的方式设置参数
func TestFunc(t *testing.T) {
	os.Args = []string{""}
	var x = conf.New()
	flag := conf.NewFlag(x)
	s := &TestFuncStruct{}
	x.RegisterConfWithName("t", s)
	x.RegisterSource(flag)
	x.Parse()
	err := x.Set("t_name", "func-name")
	assert.Nil(t, err)
	x.PrintResult()
	assert.Equal(t, "func-name", s.Name)
}

type TestSingletonStruct struct {
	String  string `conf:"string"`
	String2 string `conf:"string2"`
}

// 测试单例模式是否正常
func TestSingleton(t *testing.T) {
	os.Args = []string{"", "-t_string=CCC"}
	s := &TestSingletonStruct{}
	conf.RegisterConfWithName("t", s)
	flag := conf.NewFlag(conf.GetConf())
	conf.RegisterSource(flag)
	conf.Parse()
	err := conf.Set("t_string2", "vm50")
	assert.Nil(t, err)
	conf.PrintResult()
	assert.Equal(t, "CCC", s.String)
	assert.Equal(t, "vm50", s.String2)
}
