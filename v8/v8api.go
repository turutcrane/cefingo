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

func (v Value) AddEventListener(e EventType, h func(event *cefingo.CV8valueT) error) (err error) {

	f := v.v8value.GetValueBykey("addEventListener")
	defer cefingo.BaseRelease(f)
	cefingo.Logf("L51: addEventListener is function? :%t", f.IsFunction())

	eHander := cefingo.AllocCV8handlerT(&eventHandler{h})
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

type eventHandler struct {
	f func(event *cefingo.CV8valueT) error
}

func (eh *eventHandler) Execute(self *cefingo.CV8handlerT,
	name string,
	object *cefingo.CV8valueT,
	argumentsCount int,
	arguments []*cefingo.CV8valueT,
	retval **cefingo.CV8valueT,
	exception *cefingo.CStringT,
) (sts bool) {
	if argumentsCount == 0 {
		err := errors.Errorf("%s: No Arguments", name)
		cefingo.Logf("%+v", err)
		return false
	}
	err := eh.f(arguments[0])
	if err == nil {
		sts = true
	} else {
		cefingo.Logf("%s Not Handled %v", name, err)
	}
	return sts
}

func (v Value) HasValueBykey(key string) bool {
	return v.v8value.HasValueBykey(key)
}

func (v Value) SetValueBykey(key string, value Value) (err error) {
	if !v.v8value.SetValueBykey(key, value.v8value) {
		err = errors.Errorf("Ser value Error key:%s", key)
	}
	return err
}

func CreateString(s string) Value {
	return Value{cefingo.V8valueCreateString(s)}
}

func (c *Context) EvalString(code string) (v Value, err error) {
	var v8v *cefingo.CV8valueT
	var e *cefingo.CV8exceptionT
	if c.V8context.EvalString(code, &v8v, &e) {
		v = Value{v8v}
	} else {
		err = errors.Errorf("Eval String Error :%s", code)
	}
	return v, err
}
