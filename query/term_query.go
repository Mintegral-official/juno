package query

import (
	"fmt"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"strings"
)

type TermQuery struct {
	iterator datastruct.Iterator
	debugs   *debug.Debug
}

func NewTermQuery(iter datastruct.Iterator) (tq *TermQuery) {
	if iter == nil {
		return nil
	}
	return &TermQuery{
		iterator: iter,
	}
}

func (tq *TermQuery) Next() {
	if tq == nil || tq.iterator == nil {
		return
	}
	tq.iterator.Next()
}

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	element := tq.iterator.GetGE(id)
	if element == nil {
		if tq.debugs != nil {
			tq.debugs.AddDebugMsg(fmt.Sprintf("docId: %d, reason: %v", id, helpers.ElementNotfound))
		}
		return 0, helpers.ElementNotfound
	}
	return element.Key(), nil
}

func (tq *TermQuery) Current() (document.DocId, error) {
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}
	element := tq.iterator.Current()
	if element == nil {
		return 0, helpers.NoMoreData
	}
	return element.Key(), nil
}

func (tq *TermQuery) DebugInfo() *debug.Debug {
	if tq != nil && tq.debugs != nil {
		tq.debugs.FieldName = tq.iterator.(*datastruct.SkipListIterator).FieldName
		return tq.debugs
	}
	return nil
}

func (tq *TermQuery) Marshal() map[string]interface{} {
	if tq == nil {
		return map[string]interface{}{}
	}
	res := make(map[string]interface{}, 1)
	fields := strings.Split(tq.iterator.(*datastruct.SkipListIterator).FieldName, index.SEP)
	res["="] = []string{fields[0], fields[1]}
	return res
}

func (tq *TermQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}) Query {
	v, ok := res["="]
	if !ok {
		return nil
	}
	return NewTermQuery(idx.GetInvertedIndex().Iterator(fmt.Sprint(v.([]string)[0]), fmt.Sprint(v.([]string)[1])))
}

func (tq *TermQuery) SetDebug(level int) {
	if tq == nil {
		return
	}
	if tq.debugs == nil {
		tq.debugs = debug.NewDebug(level, "TermQuery")
	}
}
