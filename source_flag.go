package conf

import (
	"errors"
	"fmt"
	"os"
)

type Flag struct {
	*kv[interface{}]
	args []string
	conf *X
}

func NewFlag(conf *X) *Flag {
	return &Flag{
		kv:   newKV[interface{}](),
		conf: conf,
	}
}

func (f *Flag) Parse() {
	f.args = os.Args[1:]
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		f.conf.log.Errorf("parse flag err:%s", err)
	}
}

// 修改自flag标准库
func (f *Flag) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, errors.New(fmt.Sprintf("bad flag syntax: %s", s))
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	r, has := f.conf.kv.Get(name)
	if !has {
		// 没有类型无法解析
		return false, errors.New(fmt.Sprintf("flag provided but not defined: -%s", name))
	}

	if _, ok := r.(*Bool); ok { // special case: doesn't need an arg
		if !hasValue {
			value = "true"
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, errors.New(fmt.Sprintf("flag needs an argument: -%s", name))
		}
	}

	f.Set(name, value)
	return true, nil
}
