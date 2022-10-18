package v8

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/turutcrane/cefingo/capi"
)

type Context struct {
	V8context *capi.CV8contextT
	Global    Value
	Document  Value
}

type Value struct {
	v8v *capi.CV8valueT
}
type Function Value

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

func NewValue(v8v *capi.CV8valueT) Value {
	return Value{v8v: v8v}
}

func (v Value) CapiValue() *capi.CV8valueT {
	return v.v8v
}

func NewObject() Value {
	o := capi.V8valueCreateObject(nil, nil)
	return NewValue(o)
}

func NewString(s string) Value {
	return NewValue(capi.V8valueCreateString(s))
}

func GetCurrentContext() (c Context, err error) {
	v8c := capi.V8contextGetCurrentContext()
	g := NewValue(v8c.GetGlobal())
	d, err := g.GetValueBykey("document")
	if err == nil {
		c = Context{
			V8context: v8c,
			Global:    g,
			Document:  d,
		}
	}
	return c, err
}

func GetV8Context(v8c *capi.CV8contextT) (c Context, err error) {
	g := NewValue(v8c.GetGlobal())
	d, err := g.GetValueBykey("document")
	if err == nil {
		c = Context{
			V8context: v8c,
			Global:    g,
			Document:  d,
		}
	}
	return c, err
}

func (c Context) IsValid() bool {
	return c.V8context != nil && c.V8context.IsValid()
}

func (c Context) IsSame(c1 Context) bool {
	return c.IsValid() && c1.IsValid() && c.V8context.IsSame(c1.V8context)
}

func (c Context) GetBrowser() *capi.CBrowserT {
	return c.V8context.GetBrowser()
}

func (c Context) GetFrame() *capi.CFrameT {
	return c.V8context.GetFrame()
}

func (c Context) GetElementById(id string) (value Value, err error) {
	val, err := c.Document.GetValueBykey("getElementById")
	if err != nil {
		return value, err
	}
	// capi.Logf("L42: getElementById is function? :%t", f.IsFunction())

	sid := capi.V8valueCreateString(id)

	f := Function(val)
	args := []Value{NewValue(sid)}
	v8v, err := f.ExecuteFunction(c.Document, args)
	if err != nil {
		capi.Logf("L78:x %+v", err)
		return value, err
	}

	if !v8v.IsValid() || !v8v.IsObject() {
		capi.Logf("L55: Id:%s can not get valid value", id)
		err = fmt.Errorf("Id:%s can not get valid value", id)
		v8v = Value{}
	}
	return v8v, err
}

func (c Context) GetElementsByClassName(cls string) (elements Value, err error) {
	v1 := CreateString(cls)
	args := []Value{v1}

	elements, err = c.Document.Call("getElementsByClassName", args)
	if err != nil {
		capi.Logf("L77: %v", err)
	}
	return elements, err
}

func EnterContext(c *capi.CV8contextT) (ctx Context, err error) {
	runtime.LockOSThread() // UnlockOSThread at ExitContext()
	if c.Enter() {
		return GetCurrentContext()
	}
	return Context{}, fmt.Errorf("E105: Enter Error")
}

func ExitContext(c *capi.CV8contextT) error {
	current := capi.V8contextGetCurrentContext()
	if c.IsSame(current) {
		runtime.UnlockOSThread()
		if c.Exit() {
			return nil
		}
		return fmt.Errorf("E117: Context Exit Error")
	}
	return fmt.Errorf("E119: %v: Context is current %v", c, current)
}

func (v Value) IsNil() bool {
	return v.v8v == nil
}

func (v Value) AddEventListener(e EventType, h capi.CV8handlerTExecuteHandler) (err error) {

	f, err := v.GetValueBykey("addEventListener")
	if err != nil {
		return fmt.Errorf("E109: %w", err)
	}
	if !f.IsFunction() {
		return fmt.Errorf("L112: addEventListener is not function?")
	}

	eHander := capi.NewCV8handlerT(h)

	eType := capi.V8valueCreateString(string(e))

	eFunc := capi.V8valueCreateFunction("eh", eHander.Pass())

	args := []Value{NewValue(eType), NewValue(eFunc)}
	_, err = Function(f).ExecuteFunction(v, args)

	if err != nil {
		capi.Logf("T125:x %+v", err)
	}
	return err
}

type EventHandlerFunc func(object Value, event Value) error

var _ capi.CV8handlerTExecuteHandler = EventHandlerFunc(func(o, e Value) error { return nil })

func (f EventHandlerFunc) Execute(self *capi.CV8handlerT,
	name string,
	object *capi.CV8valueT,
	arguments []*capi.CV8valueT,
) (sts bool, retval *capi.CV8valueT, exception string) {
	if len(arguments) == 0 {
		err := fmt.Errorf("%s: No Arguments", name)
		capi.Logf("%+v", err)
		return false, nil, ""
	}
	err := f(Value{object}, Value{arguments[0]})
	if err == nil {
		sts = true
	} else {
		exception = fmt.Sprintf("L134: %s Not Handled %v", name, err)
		capi.Logf(exception)
	}
	return sts, retval, exception
}

func NewFunction(name string, f capi.CV8handlerTExecuteHandler) Value {
	h := capi.NewCV8handlerT(f)
	v8f := capi.V8valueCreateFunction(name, h.Pass())
	return NewValue(v8f)
}

func (f Function) ExecuteFunction(object Value, args []Value) (val Value, err error) {

	if !f.v8v.IsFunction() {
		cause := fmt.Errorf("Object is Not Function")
		return Value{}, cause
	}

	capiArgs := make([]*capi.CV8valueT, len(args))
	for i, _ := range capiArgs {
		capiArgs[i] = args[i].v8v
	}

	v8vf := f.v8v
	ret := v8vf.ExecuteFunction(object.v8v, capiArgs)
	name := v8vf.GetFunctionName()
	if ret == nil {
		if v8vf.HasException() {
			e := v8vf.GetException()
			m := e.GetMessage()
			err = fmt.Errorf("E172: %s returns NULL and %s has Exception: %s", name, name, m)
		} else if object.v8v != nil && object.HasException() {
			e := object.GetException()
			m := e.GetMessage()
			err = fmt.Errorf("E176: %s returns NULL and (this) has Exception: %s", name, m)
		} else {
			err = fmt.Errorf("E178: %s returns NULL", name)
		}
	} else if ret.IsValid() {
		val = NewValue(ret)
	} else {
		err = fmt.Errorf("E189: %s return value is not valid", name)
	}
	return val, err
}

func (f Function) ExecuteFunctionWithContext(context Context, object Value, args []Value) (val Value, err error) {

	if !f.v8v.IsFunction() {
		cause := fmt.Errorf("Object is Not Function")
		return Value{}, cause
	}

	capiArgs := make([]*capi.CV8valueT, len(args))
	for i, _ := range capiArgs {
		capiArgs[i] = args[i].v8v
	}

	v8vf := f.v8v
	ret := v8vf.ExecuteFunctionWithContext(context.V8context, object.v8v, capiArgs)
	name := v8vf.GetFunctionName()
	if ret == nil {
		if v8vf.HasException() {
			e := v8vf.GetException()
			m := e.GetMessage()
			err = fmt.Errorf("E172: %s returns NULL and %s has Exception: %s", name, name, m)
		} else if object.v8v != nil && object.HasException() {
			e := object.GetException()
			m := e.GetMessage()
			err = fmt.Errorf("E176: %s returns NULL and (this) has Exception: %s", name, m)
		} else {
			err = fmt.Errorf("E178: %s returns NULL", name)
		}
	} else if ret.IsValid() {
		val = NewValue(ret)
	} else {
		err = fmt.Errorf("E189: %s return value is not valid", name)
	}
	return val, err
}

type HandlerFunction func(this Value, args []Value) (v Value, err error)

var _ capi.CV8handlerTExecuteHandler = HandlerFunction(func(v Value, args []Value) (Value, error) { return Value{}, nil })

func (f HandlerFunction) Execute(
	self *capi.CV8handlerT,
	name string,
	thisObject *capi.CV8valueT,
	arguments []*capi.CV8valueT,
) (sts bool, retval *capi.CV8valueT, exception string) {
	sts = true // this function is handled

	args := []Value{}
	for _, a := range arguments {
		args = append(args, NewValue(a))
	}
	v, err := f(Value{thisObject}, args)
	if err == nil {
		retval = v.v8v
	} else {
		capi.Logf("L134: %s Not Handled %v", name, err)
		exception = err.Error()
	}
	return sts, retval, exception
}

func (v Value) IsValid() bool {
	if v.v8v != nil {
		return v.v8v.IsValid()
	}
	return false
}

func (v Value) IsUndefined() bool {
	if v.IsValid() {
		return v.v8v.IsUndefined()
	} else {
		return false
	}
}

func (v Value) IsNull() bool {
	if v.IsValid() {
		return v.v8v.IsNull()
	}
	return false
}

func (v Value) IsBool() bool {
	if v.IsValid() {
		return v.v8v.IsBool()
	}
	return false
}

func (v Value) IsInt() bool {
	if v.IsValid() {
		return v.v8v.IsInt()
	}
	return false
}

func (v Value) IsUnt() bool {
	if v.IsValid() {
		return v.v8v.IsUint()
	}
	return false
}

func (v Value) IsDouble() bool {
	if v.IsValid() {
		return v.v8v.IsUint()
	}
	return false
}

func (v Value) IsDate() bool {
	if v.IsValid() {
		return v.v8v.IsUint()
	}
	return false
}

func (v Value) IsString() bool {
	if v.IsValid() {
		return v.v8v.IsString()
	}
	return false
}

func (v Value) IsObject() bool {
	if v.IsValid() {
		return v.v8v.IsObject()
	}
	return false
}

func (v Value) IsFunction() bool {
	if v.IsValid() {
		return v.v8v.IsFunction()
	}
	return false
}

func (v Value) IsArray() bool {
	if v.IsValid() {
		return v.v8v.IsArray()
	}
	return false
}

func (v Value) IsArrayBuffer() bool {
	if v.IsValid() {
		return v.v8v.IsArrayBuffer()
	}
	return false
}

func (v Value) IsSame(v1 Value) bool {
	if v.IsValid() && v1.IsValid() {
		return v.v8v.IsSame(v1.v8v)
	}
	return false
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

func (v Value) GetDateValue() capi.CBasetimeT {
	return v.v8v.GetDateValue()
}

func (v Value) HasException() bool {
	return v.v8v.HasException()
}

func (v Value) GetException() *capi.CV8exceptionT {
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
		err = fmt.Errorf("Delete value Error key:%s", key)
	}
	return err
}

func (v Value) DeleteValueByindex(index int, value Value) (err error) {
	if !v.v8v.DeleteValueByindex(index) {
		err = fmt.Errorf("Delete value Error index:%d", index)
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
			err = fmt.Errorf("E360: nil returned: %s", key)
		}
	}
	return rv, err
}

func (v Value) GetValueByindex(index int) (rv Value, err error) {
	val := v.v8v.GetValueByindex(index)
	if val != nil {
		rv = Value{val}
	} else {
		if v.HasException() {
			exc := v.GetException()
			err = errors.New(exc.GetMessage())
		} else {
			err = errors.New("E375: nil returned")
		}
	}
	return rv, err
}

func (v Value) SetValueBykey(key string, value Value) (err error) {
	if !v.v8v.SetValueBykey(key, value.v8v, capi.V8PropertyAttributeNone) {
		err = fmt.Errorf("Set value Error key:%s", key)
	}
	return err
}

func (v Value) SetValueBykeyWithAttribute(key string, value Value, attribute capi.CV8PropertyattributeT) (err error) {
	if !v.v8v.SetValueBykey(key, value.v8v, attribute) {
		err = fmt.Errorf("Set value Error key:%s", key)
	}
	return err
}

func (v Value) SetValueByindex(index int, value Value) (err error) {
	if !v.v8v.SetValueByindex(index, value.v8v) {
		err = fmt.Errorf("Set value Error key:%d", index)
	}
	return err
}

func (v Value) Call(name string, args []Value) (r Value, e error) {
	f, err := v.GetValueBykey(name)
	if err != nil {
		return Value{}, err
	}
	return f.ExecuteFunction(v, args)
}

func (f Value) ExecuteFunction(this Value, args []Value) (r Value, e error) {
	capi.Logf("T340:")
	var rv Value

	if f.IsFunction() {
		rv, e = Function(f).ExecuteFunction(this, args)
		if e != nil {
			capi.Logf("T347:x %v", e)
		}
	} else {
		e = fmt.Errorf("E318: <%v> is not function", f)
		capi.Logf("T350: %v", e)
	}
	return rv, e
}

func (f Value) ExecuteFunctionWithContext(context Context, this Value, args []Value) (r Value, e error) {
	capi.Logf("T340:")
	var rv Value

	if f.IsFunction() {
		rv, e = Function(f).ExecuteFunctionWithContext(context, this, args)
		if e != nil {
			capi.Logf("T347:x %v", e)
		}
	} else {
		e = fmt.Errorf("E318: <%v> is not function", f)
		capi.Logf("T350: %v", e)
	}
	return rv, e
}

func CreateInt(i int) Value {
	return Value{capi.V8valueCreateInt(int32(i))}
}

func CreateString(s string) Value {
	return Value{capi.V8valueCreateString(s)}
}

func CreateStringFromByteArray(b []byte) Value {
	return Value{capi.V8valueCreateStringFromByteArray(b)}
}

func CreateBool(b bool) (v Value) {
	if b {
		v = Value{capi.V8valueCreateBool(1)}
	} else {
		v = Value{capi.V8valueCreateBool(0)}
	}
	return v
}

func (c Context) Eval(code string) (v Value, err error) {
	var ok bool
	var v8v *capi.CV8valueT
	var e *capi.CV8exceptionT
	if ok, v8v, e = c.V8context.Eval(code, "", 0); ok {
		v = Value{v8v}
	} else {
		err = fmt.Errorf("Eval Error<%s> %s", code, e.GetMessage())
	}
	return v, err
}

func (c Context) Alertf(message string, v ...interface{}) (err error) {

	f, err := c.Global.GetValueBykey("alert")
	if err != nil {
		return err
	}

	msg := capi.V8valueCreateString(fmt.Sprintf(message, v...))

	args := []Value{NewValue(msg)}
	_, err = Function(f).ExecuteFunction(c.Global, args)

	if err != nil {
		capi.Logf("T427:x %+v", err)
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

func (c Context) NewArray(elems ...Value) Value {
	arrayClass, err := c.Global.GetValueBykey("Array")
	if err != nil {
		capi.Panicf("E412: No Array")
	}
	if !arrayClass.IsFunction() {
		capi.Panicf("E415: Array is not function")
	}
	a := NewObject()
	a1, err := arrayClass.ExecuteFunction(a, elems)
	if err != nil {
		capi.Panicf("E420: %v", err)
	}
	return a1
}
