package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type TermQuery struct {
	iterator datastruct.Iterator
	debugs   *debug.Debugs
}

func NewTermQuery(iter datastruct.Iterator) *TermQuery {
	tq := &TermQuery{
		debugs: debug.NewDebugs(debug.NewDebug("TermQuery")),
	}
	if iter == nil {
		tq.debugs.DebugInfo.AddDebugMsg("the iterator is nil")
		return tq
	}
	tq.iterator = iter
	return tq
}

func (tq *TermQuery) Next() (document.DocId, error) {
	tq.debugs.NextNum++
	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	tq.iterator.Next()
	if element := tq.iterator.Current(); element != nil {
		v, ok := element.(*datastruct.Element)
		if !ok || v == nil || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	tq.debugs.GetNum++
	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	if v := tq.iterator.GetGE(id); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) Current() (document.DocId, error) {
	tq.debugs.CurNum++
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}
	if v := tq.iterator.Current(); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) DebugInfo() *debug.Debug {
	tq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("next has been called: %d", tq.debugs.NextNum))
	tq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("get has been called: %d", tq.debugs.GetNum))
	tq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("current has been called: %d", tq.debugs.CurNum))
	return tq.debugs.DebugInfo
}
