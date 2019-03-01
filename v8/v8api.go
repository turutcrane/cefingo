package v8

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/turutcrane/cefingo"
)

type Context struct {
	V8context *cefingo.CV8contextT
	Global    Value
	Document  Value
}

type Value struct {
	v8v *cefingo.CV8valueT
}

type EventType string

const (
	EventCancel EventType = "cancel"
	EventClick  EventType = "click"
	EventClose  EventType = "close"
	EventSubmit EventType = "submit"
)

func (v Value) HasOneRef() bool {
	return v.v8v.HasOneRef()
}

func CreateValue(v8v *cefingo.CV8valueT) Value {
	return Value{v8v: v8v}
}

func GetContext() (c *Context, err error) {
	v8c := cefingo.V8contextGetEnterdContext()
	g := CreateValue(v8c.GetGlobal())
	d, err := g.GetValueBykey("document")
	if err == nil {
		c = &Context{
			V8context: v8c,
			Global:    g,
			Document:  d,
		}
	}
	return c, err
}

func (c *Context) GetElementById(id string) (value Value, err error) {
	f, err := c.Document.GetValueBykey("getElementById")
	if err != nil {
		return Value{nil}, err
	}
	// cefingo.Logf("L42: getElementById is function? :%t", f.IsFunction())

	sid := cefingo.V8valueCreateString(id)

	args := []*cefingo.CV8valueT{sid}
	v8v, err := f.v8v.ExecuteFunction(c.Document.v8v, 1, args)
	if err != nil {
		cefingo.Logf("L36:x %+v", err)
		return Value{nil}, err
	}

	if !v8v.IsValid() || !v8v.IsObject() {
		cefingo.Logf("L55: Id:%s can not get valid value", id)
		err = fmt.Errorf("Id:%s can not get valid value", id)
	}
	return Value{v8v}, err
}

func (c *Context) GetElementsByClassName(cls string) (elements Value, err error) {
	v1 := CreateString(cls)
	args := []Value{v1}

	elements, err = c.Document.Call("getElementsByClassName", args)
	if err != nil {
		cefingo.Logf("L77: %v", err)
	}
	return elements, err
}

func (v Value) IsNil() bool {
	return v.v8v == nil
}

func (v Value) AddEventListener(e EventType, h cefingo.V8handler) (err error) {

	f := v.v8v.GetValueBykey("addEventListener")
	cefingo.Logf("L51: addEventListener is function? :%t", f.IsFunction())

	eHander := cefingo.AllocCV8handlerT(h)

	eType := cefingo.V8valueCreateString(string(e))

	eFunc := cefingo.V8valueCreateFunction("eh", eHander)

	args := []*cefingo.CV8valueT{eType, eFunc}
	_, err = f.ExecuteFunction(v.v8v, 2, args)

	if err != nil {
		cefingo.Logf("L36:x %+v", err)
	}
	return err
}

type EventHandlerFunc func(object Value, event Value) error

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
	err := f(Value{object}, Value{arguments[0]})
	if err == nil {
		sts = true
	} else {
		cefingo.Logf("L134: %s Not Handled %v", name, err)
	}
	return sts
}

func (v Value) IsValid() bool {
	return v.v8v.IsValid()
}

func (v Value) IsUndefined() bool {
	return v.v8v.IsUndefined()
}

func (v Value) IsNull() bool {
	return v.v8v.IsNull()
}

func (v Value) IsBool() bool {
	return v.v8v.IsBool()
}

func (v Value) IsInt() bool {
	return v.v8v.IsInt()
}

func (v Value) IsUnt() bool {
	return v.v8v.IsUint()
}

func (v Value) IsDouble() bool {
	return v.v8v.IsUint()
}

func (v Value) IsDate() bool {
	return v.v8v.IsUint()
}

func (v Value) IsString() bool {
	return v.v8v.IsString()
}

func (v Value) IsObject() bool {
	return v.v8v.IsObject()
}

func (v Value) IsFunction() bool {
	return v.v8v.IsFunction()
}

func (v Value) IsArray() bool {
	return v.v8v.IsArray()
}

func (v Value) IsArrayBuffer() bool {
	return v.v8v.IsArrayBuffer()
}

func (v Value) IsSame(v1 Value) bool {
	return v.v8v.IsSame(v1.v8v)
}

func (v Value) GetBoolValue() bool {
	return v.v8v.GetBoolValue()
}

func (v Value) GetIntValue() int {
	return int(v.v8v.GetIntValue())
}

func (v Value) GetUintValue() uint {
	return uint(v.v8v.GetUintValue())
}

func (v Value) GetDoubleValue() float64 {
	return v.v8v.GetDoubleValue()
}

func (v Value) GetStringValue() string {
	return v.v8v.GetStringValue()
}

func (v Value) GetDateValue() cefingo.CTimeT {
	return v.v8v.GetDateValue()
}

func (v Value) HasException() bool {
	return v.v8v.HasException()
}

func (v Value) GetException() *cefingo.CV8exceptionT {
	return v.v8v.GetException()
}

func (v Value) ClearException() bool {
	return v.v8v.ClearException()
}

func (v Value) HasValueBykey(key string) bool {
	return v.v8v.HasValueBykey(key)
}

func (v Value) DeleteValueBykey(key string, value Value) (err error) {
	if !v.v8v.DeleteValueBykey(key) {
		err = errors.Errorf("Delete value Error key:%s", key)
	}
	return err
}

func (v Value) GetValueBykey(key string) (rv Value, err error) {
	val := v.v8v.GetValueBykey(key)
	if val != nil {
		rv = Value{val}
	} else {
		if v.HasException() {
			exc := v.GetException()
			err = errors.New(exc.GetMessage())
		} else {
			err = errors.New("E262: nil returned")
		}
	}
	return rv, err
}

func (v Value) SetValueBykey(key string, value Value) (err error) {
	if !v.v8v.SetValueBykey(key, value.v8v) {
		err = errors.Errorf("Set value Error key:%s", key)
	}
	return err
}

func (v Value) Call(name string, args []Value) (r Value, e error) {
	v8args := make([]*cefingo.CV8valueT, len(args))
	for i, av := range args {
		v8args[i] = av.v8v
	}
	f, err := v.GetValueBykey(name)
	if err != nil {
		return Value{}, err
	}
	var rv *cefingo.CV8valueT
	if f.IsFunction() {
		rv, e = f.v8v.ExecuteFunction(v.v8v, len(args), v8args)
	} else {
		e = errors.Errorf("<%s> is not function", name)
	}
	return Value{rv}, e
}

func CreateInt(i int) Value {
	return Value{cefingo.V8valueCreateInt(i)}
}

func CreateString(s string) Value {
	return Value{cefingo.V8valueCreateString(s)}
}

func CreateStringFromByteArray(b []byte) Value {
	return Value{cefingo.V8valueCreateStringFromByteArray(b)}
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

func (c *Context) Alertf(message string, v ...interface{}) (err error) {

	f, err := c.Global.GetValueBykey("alert")
	if err != nil {
		return err
	}

	msg := cefingo.V8valueCreateString(fmt.Sprintf(message, v...))

	args := []*cefingo.CV8valueT{msg}
	_, err = f.v8v.ExecuteFunction(c.Global.v8v, 1, args)

	if err != nil {
		cefingo.Logf("L36:x %+v", err)
	}
	return err
}

func (v Value) ToString() (s string, err error) {
	str, e := v.Call("toString", []Value{})
	if e == nil {
		s = str.GetStringValue()
	}
	return s, e
}
