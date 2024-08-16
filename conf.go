package conf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type X struct {
	structs []*service
	sources []Source
	argTree *argTree
	kv      *kv[Arg]
	handler ConfigResultHandler
	result  []ConfigResult
}

type argTree struct {
	key   string
	value string
	child []*argTree
}

func newArgTree(key string, defValue string) *argTree {
	return &argTree{
		key:   key,
		value: defValue,
		child: make([]*argTree, 0),
	}
}

func (a *argTree) AppendChild(child *argTree) {
	a.child = append(a.child, child)
}

func New(bfs ...BuildFunc) *X {
	ret := &X{
		kv:      newKV[Arg](),
		argTree: &argTree{},
		handler: resultHandler,
	}
	for _, bf := range bfs {
		bf(ret)
	}
	return ret
}

type BuildFunc func(x *X)

func WithResultHandler(l ConfigResultHandler) BuildFunc {
	return func(x *X) {
		x.handler = l
	}
}

var (
	ErrRegisterConfigNotPtr = errors.New("register conf is not ptr")
	ErrArgSetValue          = errors.New("set value err")
	ErrFieldTypeNotSupport  = errors.New("field type not support")
	ErrArgSetDefaultValue   = errors.New("set default value err")
)

func (x *X) RegisterConf(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Ptr {
		x.handler(NewParseResultError(ErrRegisterConfigNotPtr))
		return
	}
	x.structs = append(x.structs, &service{
		Conf: f,
	})
}

func (x *X) RegisterConfWithName(name string, f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Ptr {
		x.handler(NewParseResultError(ErrRegisterConfigNotPtr))
		return
	}
	x.structs = append(x.structs, &service{
		Conf: f,
		Name: name,
	})
}

func (x *X) RegisterSource(s Source) {
	x.sources = append(x.sources, s)
}

func (x *X) Parse() {
	// 处理所有注册的结构体 创建对应的参数列表
	for _, model := range x.structs {
		x.parseStruct(model)
	}
	// 处理所有注册的配置源
	for _, source := range x.sources {
		source.Parse()
		// 将配置源中的配置参数设置到对应的参数列表中
		source.Range(func(key string, value interface{}) bool {
			arg, has := x.kv.Get(key)
			// 如果配置源中的配置参数在参数列表中不存在，那么就忽略
			if !has {
				return false
			}
			// 如果参数已经被优先级更高的配置源设置过，那么就忽略
			if arg.HasSet() {
				return false
			}
			// 将配置源中的配置参数设置到对应的参数列表中
			err := arg.SetValue(value)
			if err != nil {
				x.handler(NewParseResultError(
					ErrArgSetValue,
					errors.New(fmt.Sprintf("arg %s SetValue %v", key, value)),
					err,
				))
				return false
			}
			// 如果没有报错，那么就设置参数已经被设置过的标志
			arg.Set()
			return true
		})
	}
}

type Var struct {
	Default string
	Name    string
	Desc    string
}

type service struct {
	// 配置结构体名称
	Name string
	// 配置结构体指针
	Conf interface{}
}

func (x *X) parseStruct(service *service) {
	conf := service.Conf
	// 需要确保传入的是指针
	if reflect.TypeOf(conf).Kind() == reflect.Ptr {
		confStruct := reflect.ValueOf(conf).Elem()
		// 获取结构体的名称
		confName := snakeCase(confStruct.Type().Name())
		if service.Name == "" {
			service.Name = confName
		}
		tree := newArgTree(service.Name, "")
		x.parseTag(tree, confStruct, service.Name)
		x.argTree.AppendChild(tree)
		return
	} else {
		x.handler(NewParseResultError(ErrRegisterConfigNotPtr))
		return
	}
}

func (x *X) parseTag(tree *argTree, conf reflect.Value, tags ...string) {
	t := conf.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := conf.Field(i)

		confTag := field.Tag.Get("conf")
		// 如果配置标签为-，那么就忽略
		if confTag == "-" {
			continue
		}

		confList := strings.Split(confTag, ",")
		if field.Type.Kind() == reflect.Struct {
			// 判断是否是匿名继承
			if field.Anonymous {
				x.parseTag(tree, value, tags...)
				continue
			}
			tag := snakeCase(field.Name)
			if confTag != "" && len(confList) > 0 && !strings.Contains(confList[0], "=") {
				tag = confList[0]
			}
			newTree := newArgTree(tag, "")
			x.parseTag(newTree, value, append(tags, tag)...)
			tree.AppendChild(newTree)
			continue
		}

		attr := Var{}
		for _, keyValue := range confList {
			kvList := strings.Split(keyValue, "=")
			if len(kvList) == 1 {
				attr.Name = kvList[0]
				continue
			}
			switch kvList[0] {
			case "name":
				attr.Name = kvList[1]
			case "default":
				attr.Default = kvList[1]
			case "usage":
				attr.Desc = kvList[1]
			default:
			}
		}
		if attr.Name == "" {
			attr.Name = snakeCase(field.Name)
		}

		// 唯一键
		key := strings.Join(append(tags, attr.Name), "_")

		var arg Arg
		switch field.Type.Kind() {
		// Struct 已经在上面处理过 所以这里遇到就是错误的情况
		case reflect.Struct:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:struct", field.Name, key)),
			))
			return
		case reflect.Ptr:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:ptr", field.Name, key)),
			))
			return
		// TODO: 额外处理
		case reflect.Slice:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:slice", field.Name, key)),
			))
		// TODO: 额外处理
		case reflect.Map:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:map", field.Name, key)),
			))
		case reflect.Interface:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:interface", field.Name, key)),
			))
		case reflect.Complex64, reflect.Complex128:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type:complex", field.Name, key)),
			))
		case reflect.String:
			arg = NewString(&value)
		case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
			arg = NewInt(&value)
		case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8:
			arg = NewUint(&value)
		case reflect.Float64, reflect.Float32:
			arg = NewFloat(&value)
		case reflect.Bool:
			arg = NewBool(&value)
		default:
			x.handler(NewParseResultError(ErrFieldTypeNotSupport,
				errors.New(fmt.Sprintf("field:%s, key:%s, type: unknown", field.Name, key)),
			))
		}
		// 设置Arg默认值
		arg.SetDefaultValue(attr.Default)
		if attr.Default != "" {
			err := arg.SetValue(attr.Default)
			if err != nil {
				x.handler(NewParseResultError(ErrArgSetDefaultValue,
					errors.New(fmt.Sprintf("key:%s default:%v", key, attr.Default)),
					err,
				))
			}
		}
		// 设置Arg描述
		arg.SetDescription(attr.Desc)
		// 将该Arg注册到conf的KV中
		x.kv.Set(key, arg)
		// 将该Arg注册到tree中
		tree.AppendChild(newArgTree(attr.Name, attr.Default))
	}
}

func (x *X) Get(key string) (interface{}, bool) {
	arg, has := x.kv.Get(key)
	if !has {
		return nil, false
	}
	return arg.GetValue(), true
}

func (x *X) Set(key string, value interface{}) error {
	arg, has := x.kv.Get(key)
	if !has {
		x.kv.Set(key, NewInterface(value))
		return nil
	}
	err := arg.SetValue(value)
	return err
}

type ConfigResult struct {
	Key     string
	Value   interface{}
	Default string
	Usage   string
}

type ParseResult struct {
	Err        error
	ErrMessage string
	configs    []ConfigResult
}

func NewParseResultError(err ...error) *ParseResult {
	if len(err) == 0 {
		return NewParseResult(nil)
	}
	if len(err) == 1 {
		return &ParseResult{
			Err:     err[0],
			configs: nil,
		}
	}
	return &ParseResult{
		Err:     errors.Join(err...),
		configs: nil,
	}
}

func NewParseResult(configs []ConfigResult) *ParseResult {
	return &ParseResult{
		Err:     nil,
		configs: configs,
	}
}

type ConfigResultHandler func(*ParseResult)

func (x *X) PrintResult() {
	// 根据 argTree 进行递归打印
	x.printArgTree(x.argTree, []string{})
	x.handler(NewParseResult(x.result))
}

func (x *X) printArgTree(tree *argTree, prefix []string) {
	// 叶子节点 打印参数
	if len(tree.child) == 0 {
		key := strings.Join(append(prefix, tree.key), "_")
		arg, _ := x.kv.Get(key)
		x.result = append(x.result, ConfigResult{
			Key:     key,
			Value:   arg.GetValue(),
			Default: arg.GetDefaultValue(),
			Usage:   arg.GetDescription(),
		})
		return
	}
	// 非叶子节点 递归打印
	for _, child := range tree.child {
		nextPrefix := prefix
		if tree.key != "" {
			nextPrefix = append(prefix, tree.key)
		}
		x.printArgTree(child, nextPrefix)
	}
}
