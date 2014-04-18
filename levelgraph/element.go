package levelgraph


import (
		//"bytes"
		//"fmt"
		//"errors"
		//"os"
		//"github.com/jmhodges/levigo"
		//"github.com/lexlapax/graveldb/core"
)

type ElementType string
const (
	VertexType ElementType = "1"
	EdgeType ="2"
)

type DBElement struct {
	Db *DBGraph
	id []byte
	Elementtype ElementType
}

func (element *DBElement) Property(prop string) ([]byte) {
	return nil
}

func (element *DBElement) SetProperty(prop string, value []byte) (bool, error) {
	return false, nil	
}

func (element *DBElement) DelProperty(prop string) ([]byte) {
	return nil
}

func (element *DBElement) PropertyKeys() ([][]byte) {
	return nil
}

func (element *DBElement) Id() ([]byte) {
	return element.id
}


