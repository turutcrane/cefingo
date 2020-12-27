package router

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/turutcrane/cefingo/capi"
	"github.com/turutcrane/cefingo/cef"
	v8 "github.com/turutcrane/cefingo/v8api"
)

const (
	routerQueryMsg    = ":query"
	routerCancelMsg   = ":cancel"
	routerResponseMsg = ":response"
)

const (
	kMemberRequest    = "request"
	kMemberOnSuccess  = "onSuccess"
	kMemberOnFailure  = "onFailure"
	kMemberPersistent = "persistent"
)

type routerId int32

var routerIdGen routerId

type requestId int32

var requestIdGen requestId

var requestInfoMap sync.Map
var rendererPersistentRequsetSet sync.Map

type requestInfoKey struct {
	request requestId
}

type requestInfo struct {
	persistent      bool
	successCallback v8.Value
	failureCallback v8.Value
	router          *RendererMessageRouter
	callContext     v8.Context
	requestId       requestId
}

func requestInfoMapAdd(info *requestInfo) requestId {
	request := requestId(atomic.AddInt32((*int32)(&requestIdGen), 1))
	key := requestInfoKey{
		request,
	}
	info.requestId = request
	requestInfoMap.Store(key, info)
	if info.persistent {
		rendererPersistentRequsetSet.Store(key, true)
	}
	return request
}

func requestInfoMapGetAndDelete(browserId int, request requestId) (info *requestInfo) {
	key := requestInfoKey{
		request,
	}
	if i, got := requestInfoMap.LoadAndDelete(key); got {
		info = i.(*requestInfo)
	}
	if info != nil && info.persistent {
		rendererPersistentRequsetSet.Delete(key)
	}
	return info
}

func requestInfoMapGet(browserId int, request requestId) (info *requestInfo) {
	key := requestInfoKey{
		request,
	}
	if i, got := requestInfoMap.Load(key); got {
		info = i.(*requestInfo)
	}
	return info
}

func rendererRequestIsPersistent(browserId int, request requestId) bool {
	key := requestInfoKey{
		request,
	}
	_, loaded := requestInfoMap.Load(key)
	return loaded
}

type RendererMessageRouter struct {
	routerId                routerId
	frameContext            v8.Context
	msgPrefix               string
	queryFunctionName       string
	queryCancelFunctionName string
}

func newRendererMessageHandler(prefix string, frameContext v8.Context, queryFunctionName, queryCancelFunctionName string) *RendererMessageRouter {
	impl := &RendererMessageRouter{
		routerId:                routerId(atomic.AddInt32((*int32)(&routerIdGen), 1)),
		frameContext:            frameContext,
		msgPrefix:               prefix + ":",
		queryFunctionName:       queryFunctionName,
		queryCancelFunctionName: queryCancelFunctionName,
	}
	return impl
}

func (router *RendererMessageRouter) QueryHandler() v8.HandlerFunction {
	return v8.HandlerFunction(func(this v8.Value, args []v8.Value) (v v8.Value, err error) {
		if len(args) != 1 || !args[0].IsObject() {
			return v, fmt.Errorf("Invalid arguments; expecting a single object")
		}
		arg := args[0]
		requestVal, err := arg.GetValueBykey(kMemberRequest)
		if err != nil || !requestVal.IsString() {
			return v, fmt.Errorf("Invalid arguments; object member '%s' is required and must have type string: %w",
				kMemberRequest, err)
		}

		var successVal, failureVal, persistentVal v8.Value
		if arg.HasValueBykey(kMemberOnSuccess) {
			successVal, err = arg.GetValueBykey(kMemberOnSuccess)
			if err != nil || !successVal.IsFunction() {
				return v, fmt.Errorf("Invalid arguments; object member '%s' must have type function: %w",
					kMemberOnSuccess, err)
			}
		}

		if arg.HasValueBykey(kMemberOnFailure) {
			failureVal, err = arg.GetValueBykey(kMemberOnFailure)
			if err != nil || !failureVal.IsFunction() {
				return v, fmt.Errorf("Invalid arguments; object member '%s' must have type function: %w",
					kMemberOnFailure, err)
			}
		}

		var persistent bool
		if arg.HasValueBykey(kMemberPersistent) {
			persistentVal, err = arg.GetValueBykey(kMemberPersistent)
			if err != nil || !persistentVal.IsBool() {
				return v, fmt.Errorf("Invalid arguments; object member '%s' must have type boolean: %w",
					kMemberPersistent, err)
			}
			persistent = persistentVal.GetBoolValue()
		}

		context, err := v8.GetCurrentContext()
		if err != nil {
			capi.Panicf("T231:")
		}

		// SendQuery
		info := &requestInfo{
			persistent:      persistent,
			successCallback: successVal,
			failureCallback: failureVal,
			router:          router,
			callContext:     context,
		}
		request_id := requestInfoMapAdd(info)
		message := capi.ProcessMessageCreate(router.msgPrefix + routerQueryMsg)
		queryArgs := message.GetArgumentList()
		queryArgs.SetInt(0, int(router.routerId))
		queryArgs.SetInt(1, int(request_id))
		queryArgs.SetString(2, requestVal.GetStringValue())
		queryArgs.SetBool(3, persistent)
		context.GetFrame().SendProcessMessage(capi.PidBrowser, message)
		retVal := v8.NewValue(capi.V8valueCreateInt(int32(request_id)))
		return retVal, nil
	})
}

func (router *RendererMessageRouter) QueryCancelHandler() v8.HandlerFunction {
	return v8.HandlerFunction(func(this v8.Value, args []v8.Value) (v v8.Value, err error) {
		if len(args) != 1 || !args[0].IsInt() {
			return v, fmt.Errorf("Invalid arguments; expecting a single integer")
		}

		result := false
		request_id := requestId(args[0].GetIntValue())
		context := capi.V8contextGetCurrentContext()
		browser := context.GetBrowser()
		if request_id != 0 {
			result = sendCancel(router, browser, context.GetFrame(), request_id)
		}

		retVal := v8.CreateBool(result)
		return retVal, nil
	})
}

func sendCancel(router *RendererMessageRouter, browser *capi.CBrowserT, frame *capi.CFrameT, request requestId) (result bool) {
	browserId := browser.GetIdentifier()

	info := requestInfoMapGetAndDelete(browserId, request)
	if info != nil {
		message := capi.ProcessMessageCreate(router.msgPrefix + routerCancelMsg)
		args := message.GetArgumentList()
		args.SetInt(0, int(request))
		frame.SendProcessMessage(capi.PidBrowser, message)
		return true
	}
	return false
}

func executeSuccessCallback(browser_id int, request requestId, response string) {
	var info *requestInfo
	if rendererRequestIsPersistent(browser_id, request) {
		info = requestInfoMapGet(browser_id, request)
	} else {
		info = requestInfoMapGetAndDelete(browser_id, request)
	}

	if info == nil {
		return
	}
	if info.successCallback.IsValid() {
		args := []v8.Value{v8.CreateString(response)}
		info.successCallback.ExecuteFunctionWithContext(info.callContext, v8.Value{}, args)
	}
}

func executeFailureCallback(browser_id int, request requestId, error_code int, error_message string) {
	info := requestInfoMapGetAndDelete(browser_id, request)
	if info == nil {
		return
	}
	if info.failureCallback.IsValid() {
		args := []v8.Value{v8.CreateInt(error_code), v8.CreateString(error_message)}
		info.failureCallback.ExecuteFunctionWithContext(info.callContext, v8.Value{}, args)
	}
}

func (router *RendererMessageRouter) OnProcessMessageReceived(
	browser *capi.CBrowserT,
	frame *capi.CFrameT,
	source_process capi.CProcessIdT,
	message *capi.CProcessMessageT,
) (ret bool) {
	if !capi.CurrentlyOn(capi.TidRenderer) {
		capi.Panicf("T121:")
	}
	message_name := message.GetName()
	receivePrefix := router.msgPrefix + strconv.FormatInt(int64(router.routerId), 10)
	if strings.HasPrefix(message_name, receivePrefix) {
		msg := strings.Replace(message_name, receivePrefix, "", 1)
		args := message.GetArgumentList()
		nArgs := args.GetSize()
		request := requestId(args.GetInt(1))
		is_success := args.GetBool(2)
		switch msg {
		case routerResponseMsg:
			if nArgs <= 3 {
				capi.Panicf("T128:")
			}
			if is_success {
				if nArgs == 4 {
					response := args.GetString(3)
					cef.PostTask(capi.TidRenderer, cef.TaskFunc(func() {
						executeSuccessCallback(
							browser.GetIdentifier(),
							request,
							response,
						)
					}))
				} else {
					capi.Panicf("T136:", nArgs)
				}
			} else {
				if nArgs == 5 {
					error_code := args.GetInt(3)
					error_message := args.GetString(4)
					cef.PostTask(capi.TidRenderer, cef.TaskFunc(func() {
						executeFailureCallback(browser.GetIdentifier(), request, error_code, error_message)
					}))
				} else {
					capi.Panicf("T136:", nArgs)

				}
			}
			return true
		}
	}
	return false
}

func RendererProcessOnContextCreated(prefix string, context *capi.CV8contextT, queryFunctionName, queryCancelFunctionName string) *RendererMessageRouter {
	c, err := v8.GetCurrentContext()
	if err != nil {
		capi.Panicln("T349:", err)
	}
	// window := v8.NewValue(context.GetGlobal())
	window := c.Global
	impl := newRendererMessageHandler(prefix, c, queryFunctionName, queryCancelFunctionName)
	queryFunc := v8.NewFunction(queryFunctionName, impl.QueryHandler())
	window.SetValueBykeyWithAttribute(queryFunctionName, queryFunc,
		capi.V8PropertyAttributeReadonly|capi.V8PropertyAttributeDontenum|capi.V8PropertyAttributeDontdelete)

	cancelFunc := v8.NewFunction(queryCancelFunctionName, impl.QueryCancelHandler())
	window.SetValueBykeyWithAttribute(queryCancelFunctionName, cancelFunc,
		capi.V8PropertyAttributeReadonly|capi.V8PropertyAttributeDontenum|capi.V8PropertyAttributeDontdelete)

	return impl
}

func (router *RendererMessageRouter) OnContextReleased(
	browser *capi.CBrowserT, frame *capi.CFrameT, context *capi.CV8contextT,
) {
	if !capi.CurrentlyOn(capi.TidRenderer) {
		capi.Panicf("T121:")
	}
	keys := []requestInfoKey{}

	requestInfoMap.Range(func(key interface{}, value interface{}) bool {
		k := key.(requestInfoKey)
		info := value.(*requestInfo)
		if info.router == router && info.requestId == k.request {
			keys = append(keys, k)
		}
		return true
	})

	for _, k := range keys {
		if _, loaded := requestInfoMap.LoadAndDelete(k); loaded {
			// info := i.(*requestInfo)
			sendCancel(router, browser, context.GetFrame(), k.request)
		}
	}
}

// func sendCancelAllReauestInContext(browser *capi.CBrowserT, frame *capi.CFrameT, contextId rendererContextId) (result bool) {
// 	browserId := browser.GetIdentifier()

// 	deleted := requestInfoMapDeleteAllReauestInContext(browserId, contextId)
// 	if deleted {
// 		message := capi.ProcessMessageCreate(cancel_message_name)
// 		args := message.GetArgumentList()
// 		args.SetInt(0, int(contextId))
// 		args.SetInt(1, 0) // requestId
// 		frame.SendProcessMessage(capi.PidBrowser, message)
// 		return true
// 	}
// 	return false
// }
