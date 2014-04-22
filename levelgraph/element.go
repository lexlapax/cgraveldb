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

type ElementLevigo struct {
	db *GraphLevigo
	id []byte
	Elementtype ElementType
}

func (element *ElementLevigo) Property(prop string) ([]byte) {
	return element.db.ElementProperty(element, prop)
}

func (element *ElementLevigo) SetProperty(prop string, value []byte) (error){
	return element.db.ElementSetProperty(element, prop, value)
}

func (element *ElementLevigo) DelProperty(prop string) ([]byte) {
	return element.db.ElementDelProperty(element, prop)
}

func (element *ElementLevigo) PropertyKeys() ([]string) {
	return element.db.ElementPropertyKeys(element)
}

func (element *ElementLevigo) Id() ([]byte) {
	return element.id
}

func (element *ElementLevigo) IdAsString() (string) {
	return string(element.id[:])
}
