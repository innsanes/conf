package conf

import "fmt"

type resultLog struct {
}

func (l *resultLog) Info(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
}

type parseLog struct {
}

func (l *parseLog) Fatal(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
