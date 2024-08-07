## Prepare
通过`go get`命令下载依赖包
```shell
go get -u github.com/innsanes/conf
```
## Usage
1. 创建一个注册参数的结构体
```go
type TestStruct struct {
	String  string `conf:"string"`
	String2 string `conf:"string2"`
}
```
2. 注册参数结构体, 可以注册多个结构体, 并添加前缀进行注册
```go
s := &TestStruct{}
conf.RegisterConfWithName("t", s)
```
3. 注册配置获取的途径, 比如通过命令行获取参数或者通过yaml文件获取参数
```go
flag := conf.NewFlag(conf.GetConf())
conf.RegisterSource(flag)
```
4. 解析, 注意需要先都注册完成后再进行解析
```go
conf.Parse()
```
5. 打印结果
```go
conf.PrintResult()
```