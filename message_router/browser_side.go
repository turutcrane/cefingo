package router

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/turutcrane/cefingo/capi"
	"github.com/turutcrane/cefingo/cef"
)

type browserQueryInfo struct {
	valid      bool
	browser    capi.RefToCBrowserT
	frame      capi.RefToCFrameT
	queryId    BrowserQueryId
	msgPrefix  string
	routerId   routerId
	request    requestId
	persistent bool
	handler    BrowserQueryHandler
}

type BrowserQueryId int64

var queryIdGen BrowserQueryId

func getNextQueryId() BrowserQueryId {
	return BrowserQueryId(atomic.AddInt64((*int64)(&queryIdGen), 1))
}

type browserQueryInfoKey struct {
	browser_id int
	query_id   BrowserQueryId
}

var browserQueryInfoMap sync.Map
var browserPersistentQuerySet sync.Map

func addQueryInfoMap(browse_id int, query_id BrowserQueryId, queryInfo *browserQueryInfo) {
	key := browserQueryInfoKey{browse_id, query_id}
	browserQueryInfoMap.Store(key, queryInfo)
	if queryInfo.persistent {
		browserPersistentQuerySet.Store(key, true)
	}
}

func getAndDelteQueryInfo(browse_id int, query_id BrowserQueryId) (info *browserQueryInfo) {
	key := browserQueryInfoKey{
		browse_id,
		query_id,
	}
	if i, loaded := browserQueryInfoMap.LoadAndDelete(key); loaded {
		info = i.(*browserQueryInfo)
	}
	if info != nil && info.persistent {
		browserPersistentQuerySet.Delete(key)
	}
	return info
}

func getQueryInfo(browse_id int, query_id BrowserQueryId) (info *browserQueryInfo) {
	key := browserQueryInfoKey{
		browse_id,
		query_id,
	}
	if i, loaded := browserQueryInfoMap.Load(key); loaded {
		info = i.(*browserQueryInfo)
	}
	return info
}

func isPersistentQeury(browser_id int, query_id BrowserQueryId) bool {
	key := browserQueryInfoKey{
		browser_id,
		query_id,
	}
	_, loaded := browserPersistentQuerySet.Load(key)
	return loaded
}

type Callback interface {
	Success(response string)
	Failure(error_code int, error_message string)
	GetQueryId() BrowserQueryId
}

type BrowserQueryHandler interface {
	OnQuery(browser *capi.CBrowserT, frame *capi.CFrameT, request string, persistent bool, queryId BrowserQueryId, callback Callback) (handled bool)
	OnQueryCanceled(browser *capi.CBrowserT, frame *capi.CFrameT, queryIdd BrowserQueryId)
}

func BrowserOnProcessMessageReceived(
	handler BrowserQueryHandler,
	browser *capi.CBrowserT,
	frame *capi.CFrameT,
	msgPrefix string,
	message *capi.CProcessMessageT,
) (ret bool) {
	if !capi.CurrentlyOn(capi.TidUi) {
		log.Panicln("T380")
	}
	message_name := message.GetName()
	log.Println("T106:", message_name)
	args := message.GetArgumentList()
	browser_id := browser.GetIdentifier()
	switch message_name {
	case msgPrefix + ":" + routerQueryMsg:
		routerId := routerId(args.GetInt(0))
		request := requestId(args.GetInt(1))
		request_str := args.GetString(2)
		persistent := args.GetBool(3)

		query_id := getNextQueryId()

		// PromptHandler::OnQuery
		info := &browserQueryInfo{
			valid: true,
			// browser: capi.RefToCBrowserT{},
			// frame:  capi.RefToCFrameT{},
			queryId:    query_id,
			msgPrefix:  msgPrefix,
			routerId:   routerId,
			request:    request,
			persistent: persistent,
			handler:    handler,
		}
		info.browser.NewRefCBrowserT(browser)
		info.frame.NewRefCFrameT(frame)

		addQueryInfoMap(browser_id, query_id, info)
		if !handler.OnQuery(browser, frame, request_str, persistent, query_id, info) {
			cancelUnhandledQuery(info)
		}
		return true

	case msgPrefix + routerCancelMsg:
		routerId := routerId(args.GetInt(0))
		request := requestId(args.GetInt(1))
		cancelPendingRequest(browser_id, routerId, request)
		return true
	}

	return false
}

func (info *browserQueryInfo) Success(response string) {
	if !capi.CurrentlyOn(capi.TidUi) {
		cef.PostTask(capi.TidUi, cef.TaskFunc(func() {
			info.Success(response)
		}))
	}

	if info.valid {
		if info.persistent {
			info = getQueryInfo(info.browser.GetCBrowserT().GetIdentifier(), info.queryId)
		} else {
			info = getAndDelteQueryInfo(info.browser.GetCBrowserT().GetIdentifier(), info.queryId)
		}
		if info != nil {
			cef.PostTask(capi.TidUi, cef.TaskFunc(func() {
				sendQuerySuccess(info, response)
			}))
		}
	}
}

func sendQuerySuccess(info *browserQueryInfo, response string) {
	message := capi.ProcessMessageCreate(info.msgPrefix + strconv.FormatInt(int64(info.routerId), 10) + routerResponseMsg)
	args := message.GetArgumentList()
	args.SetInt(0, int(info.routerId))
	args.SetInt(1, int(info.request))
	args.SetBool(2, true)
	args.SetString(3, response)
	info.frame.GetCFrameT().SendProcessMessage(capi.PidRenderer, message)
	if !info.persistent {
		info.browser.UnrefCBrowserT()
		info.frame.UnrefCFrameT()
	}
}

func (info *browserQueryInfo) Failure(errorCode int, errorMessage string) {
	if !capi.CurrentlyOn(capi.TidUi) {
		cef.PostTask(capi.TidUi, cef.TaskFunc(func() {
			info.Failure(errorCode, errorMessage)
		}))
	}
	if info.valid {
		info := getAndDelteQueryInfo(info.browser.GetCBrowserT().GetIdentifier(), info.queryId)
		if info != nil {
			cef.PostTask(capi.TidUi, cef.TaskFunc(func() {
				sendQueryFailure(info,
					errorCode, errorMessage,
				)
			}))
		}
	}
}

func (info *browserQueryInfo) GetQueryId() BrowserQueryId {
	return info.queryId
}

func sendQueryFailure(info *browserQueryInfo, errorCode int, errorMessage string) {
	messageName := info.msgPrefix + strconv.FormatInt(int64(info.routerId), 10) + routerResponseMsg
	message := capi.ProcessMessageCreate(messageName)
	args := message.GetArgumentList()
	args.SetInt(0, int(info.routerId))
	args.SetInt(1, int(info.request))
	args.SetBool(2, false)
	args.SetInt(3, errorCode)
	args.SetString(4, errorMessage)
	info.frame.GetCFrameT().SendProcessMessage(capi.PidRenderer, message)
}

const (
	kCanceledErrorCode    = -1
	kCanceledErrorMessage = "The query has been canceled"
)

func cancelUnhandledQuery(info *browserQueryInfo) {
	sendQueryFailure(info, kCanceledErrorCode, kCanceledErrorMessage)
	info.browser.UnrefCBrowserT()
	info.frame.UnrefCFrameT()
}

func cancelPendingRequest(browser_id int, routerId routerId, request requestId) {
	keys := []browserQueryInfoKey{}
	browserQueryInfoMap.Range(func(key interface{}, value interface{}) bool {
		k := key.(browserQueryInfoKey)
		info := value.(browserQueryInfo)
		if k.browser_id == browser_id && info.routerId == routerId && (request == 0 || info.request == request) {
			keys = append(keys, k)
		}
		return true
	})
	for _, k := range keys {
		if i, loaded := browserQueryInfoMap.LoadAndDelete(k); loaded {
			info := i.(*browserQueryInfo)
			cancelQuery(k.query_id, info, false)
		}
	}
}

// call on OnBeforClose, OnRenderProcessTerminated, OnBeforeBrowse(if frame.IsMain() is true)
func BrowserCancelPendingForBrowser(browser *capi.CBrowserT) {
	if !capi.CurrentlyOn(capi.TidUi) {
		cef.PostTask(capi.TidUi, cef.TaskFunc(func() {
			BrowserCancelPendingForBrowser(browser)
		}))
	}
	browser_id := browser.GetIdentifier()

	keys := []browserQueryInfoKey{}
	browserQueryInfoMap.Range(func(key interface{}, value interface{}) bool {
		k := key.(browserQueryInfoKey)
		if k.browser_id == browser_id {
			keys = append(keys, k)
		}
		return true
	})
	for _, k := range keys {
		if i, loaded := browserQueryInfoMap.LoadAndDelete(k); loaded {
			info := i.(*browserQueryInfo)
			cancelQuery(k.query_id, info, false)
		}
	}
}

func cancelQuery(queryId BrowserQueryId, info *browserQueryInfo, notifyRederer bool) {
	if notifyRederer {
		sendQueryFailure(info, kCanceledErrorCode, kCanceledErrorMessage)
	}
	info.handler.OnQueryCanceled(info.browser.GetCBrowserT(), info.frame.GetCFrameT(), queryId)
	info.valid = false
	info.browser.UnrefCBrowserT()
	info.frame.UnrefCFrameT()
}
