package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrOpenFile  = errors.New("open file err")
	ErrReadFile  = errors.New("read file err")
	ErrWriteFile = errors.New("write file err")
)

func (x *X) readFile(filepath string) []byte {
	// 打开文件
	file, err := os.Open(filepath)
	if err != nil {
		x.handler(NewParseResultError(ErrOpenFile, err))
	}
	defer func() {
		_ = file.Close()
	}()
	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		x.handler(NewParseResultError(ErrReadFile, err))
	}
	return content
}

func (x *X) fileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (x *X) writeFile(filepath string, content []byte) {
	// 打开文件
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		x.handler(NewParseResultError(ErrOpenFile, err))
	}
	defer func() {
		_ = file.Close()
	}()
	// 写入文件内容
	_, err = file.Write(content)
	if err != nil {
		x.handler(NewParseResultError(ErrWriteFile, err))
	}
}

func snakeCase(str string) string {
	var (
		ret []rune
	)
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			ret = append(ret, '_')
		}
		ret = append(ret, r)
	}
	return strings.ToLower(string(ret))
}

func (x *X) panic(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
