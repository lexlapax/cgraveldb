package levelgraph


import (
		//"fmt"
		//"github.com/lexlapax/graveldb/core"
)

type ElementType string
const (
	VertexType ElementType = "1"
	EdgeType ="2"
)

type DBElement struct {
	db *DBGraph
	id []byte
	Elementtype ElementType
}

func (element *DBElement) Property(prop string) ([]byte) {
	return element.db.ElementProperty(element, prop)
}

func (element *DBElement) SetProperty(prop string, value []byte) (error){
	return element.db.ElementSetProperty(element, prop, value)
}

func (element *DBElement) DelProperty(prop string) ([]byte) {
	return element.db.ElementDelProperty(element, prop)
}

func (element *DBElement) PropertyKeys() ([]string) {
	return element.db.ElementPropertyKeys(element)
}

func (element *DBElement) Id() ([]byte) {
	return element.id
}

func (element *DBElement) IdAsString() (string) {
	return string(element.id[:])
}
