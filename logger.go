package conf

import "fmt"

type log struct {
}

func (l *log) Fatalf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}

func (l *log) Errorf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf("ERR: "+format, v...))
}

func (l *log) Infof(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
}
