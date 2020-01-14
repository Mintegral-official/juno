package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"strconv"
)

type TermQuery struct {
	iterator datastruct.Iterator
	debugs   *debug.Debugs
}

func NewTermQuery(iter datastruct.Iterator, isDebug ...int) (tq *TermQuery) {
	tq = &TermQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		tq.debugs = debug.NewDebugs(debug.NewDebug("TermQuery"))
	}
	if iter == nil {
		tq.debugs.DebugInfo.AddDebugMsg("the iterator is nil")
		return tq
	}
	tq.iterator = iter
	return tq
}

func (tq *TermQuery) Next() (document.DocId, error) {
	if tq.debugs != nil {
		tq.debugs.NextNum++
	}
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	tq.iterator.Next()
	element := tq.iterator.Current()
	if element == nil || element.Key() == 0 {
		return 0, helpers.ElementNotfound
	}
	return element.Key(), nil
}

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	if tq.debugs != nil {
		tq.debugs.GetNum++
	}
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	element := tq.iterator.GetGE(id)
	if element == nil || element.Key() == 0 {
		return 0, helpers.ElementNotfound
	}
	return element.Key(), nil
}

func (tq *TermQuery) Current() (document.DocId, error) {
	if tq.debugs != nil {
		tq.debugs.CurNum++
	}
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}
	element := tq.iterator.Current()
	if element == nil || element.Key() == 0 {
		return 0, helpers.ElementNotfound
	}
	return element.Key(), nil
}

func (tq *TermQuery) DebugInfo() *debug.Debug {
	if tq.debugs != nil {
		tq.debugs.DebugInfo.AddDebugMsg("next has been called: " + strconv.Itoa(tq.debugs.NextNum))
		tq.debugs.DebugInfo.AddDebugMsg("get has been called: " + strconv.Itoa(tq.debugs.GetNum))
		tq.debugs.DebugInfo.AddDebugMsg("current has been called: " + strconv.Itoa(tq.debugs.CurNum))
		return tq.debugs.DebugInfo
	}
	return nil
}
