package conf

type Source interface {
	Get(str string) (interface{}, bool)
	Parse()
	Range(f func(key string, value interface{}) bool)
}

type Arg interface {
	SetValue(v interface{}) error
	GetValue() interface{}
	SetDefaultValue(str string)
	GetDefaultValue() string
	SetDescription(str string)
	GetDescription() string
	HasSet() bool
	Set()
}

type Logger interface {
	Fatalf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}
