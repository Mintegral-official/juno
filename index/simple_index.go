package index

import "github.com/Mintegral-official/juno/query"

type SimpleIndex struct {
}

func (si *SimpleIndex) Add(doc *DocInfo) error {
	return nil
}

func (si *SimpleIndex) Del(doc *DocInfo) error {
	return nil
}

func (si *SimpleIndex) Update(filename string) error {
	return nil
}

func (si *SimpleIndex) Dump(filename string) error {
	return nil
}

func (si *SimpleIndex) Load(filename string) error {
	return nil
}

func (si *SimpleIndex) Search(query *query.Query) *query.SearchResult {
	return nil
}
