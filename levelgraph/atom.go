package levelgraph


import (
		//"fmt"
		//"github.com/lexlapax/graveldb/core"
)

type NodeType string
const (
	VertexType NodeType = "1"
	EdgeType ="2"
)

type AtomLevigo struct {
	db *GraphLevigo
	id []byte
	nodeType NodeType
}

func (node *AtomLevigo) Property(prop string) ([]byte) {
	return node.db.ElementProperty(node, prop)
}

func (node *AtomLevigo) SetProperty(prop string, value []byte) (error){
	return node.db.ElementSetProperty(node, prop, value)
}

func (node *AtomLevigo) DelProperty(prop string) ([]byte) {
	return node.db.ElementDelProperty(node, prop)
}

func (node *AtomLevigo) PropertyKeys() ([]string) {
	return node.db.ElementPropertyKeys(node)
}

func (node *AtomLevigo) Id() ([]byte) {
	return node.id
}

func (node *AtomLevigo) IdAsString() (string) {
	return string(node.id[:])
}
