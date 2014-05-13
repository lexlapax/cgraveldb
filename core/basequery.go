package core

import (
		"github.com/lexlapax/graveldb/util"
		"bytes"
)

type ConditionContainer struct {
	Key string
	Value []byte
	Has bool // has = true, nothas = false
}

type BaseQuery struct{
	Conditions []*ConditionContainer
	Max int
}

func newBaseQuery() *BaseQuery {
	query := &BaseQuery{}
	query.Conditions = []*ConditionContainer{}
	query.Max = util.MaxInt
	return query
}

func (query *BaseQuery) ConditionFilter(atom Atom) bool {
	atomkeyset := util.NewStringSet()
	propkeys, _ := atom.PropertyKeys()
	atomkeyset.AddArray(propkeys)

	for _, cond := range query.Conditions {
		if cond.Has {
			if ! atomkeyset.Contains(cond.Key) { return false }
			if cond.Value != nil {
				atompropval, _ := atom.Property(cond.Key)
				if bytes.Compare(cond.Value, atompropval) != 0 { return false }
			}
		} else {
			if atomkeyset.Contains(cond.Key) {
				if cond.Value == nil { return false }
				atompropval, _ := atom.Property(cond.Key)
				if bytes.Compare(cond.Value, atompropval) == 0 { return false }
			}
		}
	}
	return true
}

func (query *BaseQuery) HasKey(key string) Query {
	if key == "" { return query }
	queryc := &ConditionContainer{key,nil,true}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) NoKey(key string) Query {
	if key == "" { return query }
	queryc := &ConditionContainer{key,nil,false}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) HasKeyValue(key string, value []byte) Query {
	if key == "" { return query }
	queryc := &ConditionContainer{key,value,true}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) NoKeyValue(key string, value []byte) Query {
	if key == "" { return query }
	queryc := &ConditionContainer{key,value,false}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) Limit(limit int) Query {
	query.Max = limit
	return query
}

func (query *BaseQuery) IterEdges() <-chan Edge {
	return nil
}

func (query *BaseQuery) IterVertices() <-chan Vertex {
	return nil
}