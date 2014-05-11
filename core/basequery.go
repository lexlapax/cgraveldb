package core

import (
		"github.com/lexlapax/graveldb/util"
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

func (query *BaseQuery) HasKey(key string) *BaseQuery {
	if key == "" { return query }
	queryc := &ConditionContainer{key,nil,true}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) NoKey(key string) *BaseQuery {
	if key == "" { return query }
	queryc := &ConditionContainer{key,nil,false}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) HasKeyValue(key string, value []byte) *BaseQuery {
	if key == "" { return query }
	queryc := &ConditionContainer{key,value,true}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) NoKeyValue(key string, value []byte) *BaseQuery {
	if key == "" { return query }
	queryc := &ConditionContainer{key,value,false}
	query.Conditions = append(query.Conditions, queryc)
	return query
}

func (query *BaseQuery) Limit(limit int) *BaseQuery {
	query.Max = limit
	return query
}

