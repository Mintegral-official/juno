package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type TermQuery struct {
	iterator datastruct.Iterator
	debugs   *debug.Debug
}

func NewTermQuery(iter datastruct.Iterator, isDebug ...int) (tq *TermQuery) {
	tq = &TermQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		tq.debugs = debug.NewDebug("TermQuery")
	}
	if iter == nil {
		tq.debugs.AddDebugMsg("the iterator is nil")
		return tq
	}
	tq.iterator = iter
	return tq
}

func (tq *TermQuery) Next() (document.DocId, error) {
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	tq.iterator.Next()
	element := tq.iterator.Current()
	if element == nil {
		return 0, helpers.ElementNotfound
	}
	if tq.debugs != nil {
		tq.debugs.Node[element.Key()] = append(tq.debugs.Node[element.Key()],
			[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: found id"})
	}
	return element.Key(), nil
}

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	element := tq.iterator.GetGE(id)
	if element == nil {
		if tq.debugs != nil {
			tq.debugs.Node[id] = append(tq.debugs.Node[id],
				[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: not found"})
		}
		return 0, helpers.ElementNotfound
	}
	if tq.debugs != nil {
		if element.Key() != id {
			tq.debugs.Node[element.Key()] = append(tq.debugs.Node[element.Key()],
				[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: found id"})
			tq.debugs.Node[id] = append(tq.debugs.Node[id],
				[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: not found"})
		} else {
			tq.debugs.Node[id] = append(tq.debugs.Node[id],
				[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: found id"})
		}
	}
	return element.Key(), nil
}

func (tq *TermQuery) Current() (document.DocId, error) {
	if tq == nil || tq.iterator == nil {
		return 0, helpers.DocumentError
	}
	element := tq.iterator.Current()
	if element == nil {
		return 0, helpers.ElementNotfound
	}
	if tq.debugs != nil {
		tq.debugs.Node[element.Key()] = append(tq.debugs.Node[element.Key()],
			[]string{"field:" + tq.iterator.(*datastruct.SkipListIterator).FieldName, "reason: found id"})
	}
	return element.Key(), nil
}

func (tq *TermQuery) DebugInfo() *debug.Debug {
	if tq.debugs != nil {
		return tq.debugs
	}
	return nil
}

func (tq *TermQuery) Marshal(idx *index.Indexer) map[string]interface{} {
	invertIdx := idx.GetInvertedIndex().(*index.InvertedIndexer)
	if len(invertIdx.GetField()) == 0 || len(invertIdx.GetValue()) == 0 {
		return nil
	}
	field, value := invertIdx.GetField(), invertIdx.GetValue()
	res := make(map[string]interface{}, 1)
	res["="] = []string{field[0], value[0]}
	field = append(field[:0], field[1:]...)
	value = append(value[:0], value[1:]...)
	return res
}

func (tq *TermQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	v, ok := res["="]
	if !ok {
		return nil
	}
	return NewTermQuery(idx.GetInvertedIndex().Iterator(fmt.Sprint(v.([]string)[0]), fmt.Sprint(v.([]string)[1])), 1)
}

func (tq *TermQuery) SetDebug(isDebug ...int) {
	if len(isDebug) == 1 && isDebug[0] == 1 {
		tq.debugs = debug.NewDebug("TermQuery")
	}
}

func (tq *TermQuery) UnsetDebug() {
	tq.debugs = nil
}
