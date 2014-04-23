package mem

import (
		//"github.com/lexlapax/graveldb/core"
		"sync"
)

type NodeType string
const (
	VertexType NodeType = "1"
	EdgeType ="2"
)

type NodeMem struct {
	graph *GraphMem
	id []byte
	elementType NodeType
	props map[string][]byte 
	sync.RWMutex
}

func NewNodeMem(graph *GraphMem, id []byte, elementType NodeType) *NodeMem {
	node := new(NodeMem)
	node.id = id 
	node.graph = graph
	node.elementType = elementType
	node.props = make(map[string][]byte)
	return node
}

func (node *NodeMem) Id() []byte {
	return node.Id()
}

func (node *NodeMem) Property(prop string) ([]byte, error) {
	node.RLock()
	defer node.RUnlock()
 	if val, ok := node.props[prop]; ok {
        return val, nil
    }
	return nil, nil
}

func (node *NodeMem) SetProperty(prop string, value []byte) error {
	node.Lock()
	defer node.Unlock()
	node.props[prop] = value
	return nil
}

func (node *NodeMem) DelProperty(prop string) error {
	node.Lock()
	defer node.Unlock()
 	delete(node.props, prop)
	return nil
}

func (node *NodeMem) PropertyKeys() ([]string, error) {
	node.RLock()
	defer node.RUnlock()
	keys := []string{}
	for k := range node.props {
		keys = append(keys, k)
	}
	return keys, nil
}
