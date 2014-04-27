package levigo


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

func (node *AtomLevigo) Id() []byte {
	return node.id
}

func (node *AtomLevigo) Property(prop string) ([]byte, error) {
	return node.db.AtomProperty(node, prop)
}

func (node *AtomLevigo) SetProperty(prop string, value []byte) error {
	return node.db.AtomSetProperty(node, prop, value)
}

func (node *AtomLevigo) DelProperty(prop string) error {
	return node.db.AtomDelProperty(node, prop)
}

func (node *AtomLevigo) PropertyKeys() ([]string, error) {
	return node.db.AtomPropertyKeys(node)
}

func (node *AtomLevigo) IdAsString() string {
	return string(node.id[:])
}
