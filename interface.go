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

type ParseLogger interface {
	Fatal(format string, v ...interface{})
}

type ResultLogger interface {
	Info(format string, v ...interface{})
}
