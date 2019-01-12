package v8

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/turutcrane/cefingo"
)

type Context struct {
	V8context *cefingo.CV8contextT
	Global    *cefingo.CV8valueT
	Document  *cefingo.CV8valueT
}

type Value struct {
	v8value *cefingo.CV8valueT
}

type EventType string

const (
	EventClick EventType = "click"
)

func GetContext() *Context {
	c := cefingo.V8contextGetEnterdContext()
	g := c.GetGlobal()
	d := g.GetValueBykey("document")
	return &Context{c, g, d}
}

func ReleaseContext(c *Context) {
	cefingo.BaseRelease(c.Document)
	cefingo.BaseRelease(c.Global)
	cefingo.BaseRelease(c.V8context)
}

func (c *Context) GetElementById(id string) (value Value, err error) {
	f := c.Document.GetValueBykey("getElementById")
	defer cefingo.BaseRelease(f)
	cefingo.Logf("L42: getElementById is function? :%t", f.IsFunction())

	ids := cefingo.V8valueCreateString(id)
	defer cefingo.BaseRelease(ids)

	args := []*cefingo.CV8valueT{ids}
	v8value, err := f.ExecuteFunction(c.Document, 1, args)
	if err != nil {
		cefingo.Logf("L36:x %+v", err)
		return Value{nil}, err
	}

	if !v8value.IsValid() || !v8value.IsObject() {
		cefingo.Logf("L55: Id:%s can not get valid value", id)
		err = fmt.Errorf("Id:%s can not get valid value", id)
	}
	return Value{v8value}, err
}

func (v Value) AddRef() {
	cefingo.BaseAddRef(v.v8value)
}

func (v Value) Release() {
	cefingo.BaseRelease(v.v8value)
}

func (v Value) AddEventListener(e EventType, h cefingo.V8handler) (err error) {

	f := v.v8value.GetValueBykey("addEventListener")
	defer cefingo.BaseRelease(f)
	cefingo.Logf("L51: addEventListener is function? :%t", f.IsFunction())

	eHander := cefingo.AllocCV8handlerT(h)
	defer cefingo.BaseRelease(eHander)

	eType := cefingo.V8valueCreateString(string(e))
	defer cefingo.BaseRelease(eType)

	eFunc := cefingo.V8valueCreateFunction("eh", eHander)
	defer cefingo.BaseRelease(eFunc)

	args := []*cefingo.CV8valueT{eType, eFunc}
	_, err = f.ExecuteFunction(v.v8value, 2, args)

	if err != nil {
		cefingo.Logf("L36:x %+v", err)
	}
	return err
}

type EventHandlerFunc  func(event *cefingo.CV8valueT) error

func (f EventHandlerFunc) Execute(self *cefingo.CV8handlerT,
	name string,
	object *cefingo.CV8valueT,
	argumentsCount int,
	arguments []*cefingo.CV8valueT,
	retval **cefingo.CV8valueT,
	exception *string,
) (sts bool) {
	if argumentsCount == 0 {
		err := errors.Errorf("%s: No Arguments", name)
		cefingo.Logf("%+v", err)
		return false
	}
	err := f(arguments[0])
	if err == nil {
		sts = true
	} else {
		cefingo.Logf("%s Not Handled %v", name, err)
	}
	return sts
}

func (v Value) IsValid() bool {
	return v.v8value.IsValid()
}

func (v Value) IsUndefined() bool {
	return v.v8value.IsUndefined()
}

func (v Value) IsNull() bool {
	return v.v8value.IsNull()
}

func (v Value) IsBool() bool {
	return v.v8value.IsBool()
}

func (v Value) IsInt() bool {
	return v.v8value.IsInt()
}

func (v Value) IsUnt() bool {
	return v.v8value.IsUint()
}

func (v Value) IsDouble() bool {
	return v.v8value.IsUint()
}

func (v Value) IsDate() bool {
	return v.v8value.IsUint()
}

func (v Value) IsString() bool {
	return v.v8value.IsString()
}

func (v Value) IsObject() bool {
	return v.v8value.IsObject()
}

func (v Value) IsFunction() bool {
	return v.v8value.IsFunction()
}

func (v Value) IsArray() bool {
	return v.v8value.IsArray()
}

func (v Value) IsArrayBuffer() bool {
	return v.v8value.IsArrayBuffer()
}

func (v Value) IsSame(v1 Value) bool {
	return v.v8value.IsSame(v1.v8value)
}

func (v Value) GetBoolValue() bool {
	return v.v8value.GetBoolValue()
}

func (v Value) GetIntValue() int {
	return int(v.v8value.GetIntValue())
}

func (v Value) GetUintValue() uint {
	return uint(v.v8value.GetUintValue())
}

func (v Value) GetDoubleValue() float64 {
	return v.v8value.GetDoubleValue()
}

func (v Value) GetStringValue() string {
	return v.v8value.GetStringValue()
}

func (v Value) GetDateValue() cefingo.CTimeT {
	return v.v8value.GetDateValue()
}

func (v Value) HasExceptin() bool {
	return v.v8value.HasException()
}

func (v Value) GetExceptin() *cefingo.CV8exceptionT {
	return v.v8value.GetException()
}

func (v Value) ClearExceptin() bool {
	return v.v8value.ClearException()
}

func (v Value) HasValueBykey(key string) bool {
	return v.v8value.HasValueBykey(key)
}

func (v Value) DeleteValueBykey(key string, value Value) (err error) {
	if !v.v8value.DeleteValueBykey(key) {
		err = errors.Errorf("Delete value Error key:%s", key)
	}
	return err
}

func (v Value) GetValueBykey(key string) (rv Value) {
	return Value{v.v8value.GetValueBykey(key)}
}

func (v Value) SetValueBykey(key string, value Value) (err error) {
	if !v.v8value.SetValueBykey(key, value.v8value) {
		err = errors.Errorf("Set value Error key:%s", key)
	}
	return err
}

func CreateString(s string) Value {
	return Value{cefingo.V8valueCreateString(s)}
}

func (c *Context) Eval(code string) (v Value, err error) {
	var v8v *cefingo.CV8valueT
	var e *cefingo.CV8exceptionT
	if c.V8context.Eval(code, &v8v, &e) {
		v = Value{v8v}
	} else {

		err = errors.Errorf("Eval Error<%s> %s", code, e.GetMessage())
	}
	return v, err
}

func (c Context) Alertf(message string, v ...interface{}) (err error) {

	f := c.Global.GetValueBykey("alert")
	defer cefingo.BaseRelease(f)

	msg := cefingo.V8valueCreateString(fmt.Sprintf(message, v...))
	defer cefingo.BaseRelease(msg)

	args := []*cefingo.CV8valueT{msg}
	_, err = f.ExecuteFunction(c.Global, 1, args)

	if err != nil {
		cefingo.Logf("L36:x %+v", err)
	}
	return err
}
