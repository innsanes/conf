package conf

// 单例模式
var nx = New()

func GetConf() *X {
	return nx
}

func RegisterConf(f interface{}) {
	nx.RegisterConf(f)
}

func RegisterConfWithName(name string, f interface{}) {
	nx.RegisterConfWithName(name, f)
}

func RegisterSource(s Source) {
	nx.RegisterSource(s)
}

func Get(str string) (interface{}, bool) {
	return nx.Get(str)
}

func Set(str string, v interface{}) error {
	return nx.Set(str, v)
}

func Parse() {
	nx.Parse()
}

func PrintResult() {
	nx.PrintResult()
}
