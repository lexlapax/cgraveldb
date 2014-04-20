package levelgraph


import (
		//"bytes"
		"fmt"
		//"errors"
		//"os"
		//"github.com/jmhodges/levigo"
		//"github.com/lexlapax/graveldb/core"
)

type DBEdge struct {
	*DBElement
	subject *DBVertex
	object *DBVertex
	label string
}


func (edge *DBEdge) Label() (string) {
	return edge.label
}

func (edge *DBEdge) VertexOut() (*DBVertex) {
	return edge.subject
}

func (edge *DBEdge) VertexIn() (*DBVertex) {
	return edge.object
}

func (edge *DBEdge) String() (string) {
	str := fmt.Sprintf("<DBEdge:%v,s=%v,o=%v,l=%v@%v>",edge.IdAsString(), 
		edge.subject.IdAsString(),edge.object.IdAsString(),edge.label,edge.db)
	return str
}
