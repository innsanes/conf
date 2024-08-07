package conf

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrInvalidValue = errors.New("invalid value")

type Bool struct {
	rValue *reflect.Value
	DefValue
	Description
	Has
}

func NewBool(r *reflect.Value) *Bool {
	ret := &Bool{rValue: r}
	return ret
}

func (b *Bool) GetValue() interface{} {
	return b.rValue.Bool()
}

func (b *Bool) SetValue(str interface{}) error {
	switch v := str.(type) {
	case string:
		if str == "true" {
			b.rValue.SetBool(true)
			return nil
		} else if str == "false" {
			b.rValue.SetBool(false)
			return nil
		} else {
			return ErrInvalidValue
		}
	case bool:
		b.rValue.SetBool(v)
		return nil
	default:
		return ErrInvalidValue
	}
}

type Int struct {
	rValue *reflect.Value
	DefValue
	Description
	Has
}

func NewInt(r *reflect.Value) *Int {
	ret := &Int{rValue: r}
	return ret
}

func (i *Int) GetValue() interface{} {
	return i.rValue.Int()
}

func (i *Int) SetValue(str interface{}) error {
	switch v := str.(type) {
	case string:
		vv, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		i.rValue.SetInt(int64(vv))
	case int:
		i.rValue.SetInt(int64(v))
	case int64:
		i.rValue.SetInt(v)
	case int32:
		i.rValue.SetInt(int64(v))
	case int16:
		i.rValue.SetInt(int64(v))
	case int8:
		i.rValue.SetInt(int64(v))
	default:
		return ErrInvalidValue
	}
	return nil
}

type Uint struct {
	rValue *reflect.Value
	DefValue
	Description
	Has
}

func NewUint(r *reflect.Value) *Uint {
	ret := &Uint{rValue: r}
	return ret
}

func (u *Uint) GetValue() interface{} {
	return u.rValue.Uint()
}

func (u *Uint) SetValue(str interface{}) error {
	switch v := str.(type) {
	case string:
		vv, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		u.rValue.SetUint(uint64(vv))
	case uint:
		u.rValue.SetUint(uint64(v))
	case uint64:
		u.rValue.SetUint(v)
	case uint32:
		u.rValue.SetUint(uint64(v))
	case uint16:
		u.rValue.SetUint(uint64(v))
	case uint8:
		u.rValue.SetUint(uint64(v))
	default:
		return ErrInvalidValue
	}
	return nil
}

type String struct {
	rValue *reflect.Value
	DefValue
	Description
	Has
}

func NewString(r *reflect.Value) *String {
	ret := &String{rValue: r}
	return ret
}

func (s *String) GetValue() interface{} {
	return s.rValue.String()
}

func (s *String) SetValue(str interface{}) error {
	switch v := str.(type) {
	case string:
		s.rValue.SetString(v)
		return nil
	default:
		return ErrInvalidValue
	}
}

type Float struct {
	rValue *reflect.Value
	DefValue
	Description
	Has
}

func NewFloat(r *reflect.Value) *Float {
	ret := &Float{rValue: r}
	return ret
}

func (f *Float) GetValue() interface{} {
	return f.rValue.Float()
}

func (f *Float) SetValue(str interface{}) error {
	switch v := str.(type) {
	case string:
		vv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		f.rValue.SetFloat(vv)
	case float64:
		f.rValue.SetFloat(v)
	case float32:
		f.rValue.SetFloat(float64(v))
	default:
		return ErrInvalidValue
	}
	return nil
}

type Interface struct {
	value interface{}
	DefValue
	Description
	Has
}

func NewInterface(r interface{}) *Interface {
	ret := &Interface{value: r}
	return ret
}

func (i *Interface) GetValue() interface{} {
	return i.value
}

func (i *Interface) SetValue(str interface{}) error {
	i.value = str
	return nil
}

type Has struct {
	hasSet bool
}

func (h *Has) HasSet() bool {
	return h.hasSet
}

func (h *Has) Set() {
	h.hasSet = true
}

type DefValue struct {
	defValue string
}

func (d *DefValue) GetDefaultValue() string {
	return d.defValue
}

func (d *DefValue) SetDefaultValue(defValue string) {
	d.defValue = defValue
}

type Description struct {
	desc string
}

func (d *Description) GetDescription() string {
	return d.desc
}

func (d *Description) SetDescription(desc string) {
	d.desc = desc
}
