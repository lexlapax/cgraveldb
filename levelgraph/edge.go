package levelgraph


import (
		"fmt"
		//"github.com/lexlapax/graveldb/core"
)

type EdgeLevigo struct {
	*AtomLevigo
	subject *VertexLevigo
	object *VertexLevigo
	label string
}


func (edge *EdgeLevigo) Label() (string) {
	return edge.label
}

func (edge *EdgeLevigo) VertexOut() (*VertexLevigo) {
	return edge.subject
}

func (edge *EdgeLevigo) VertexIn() (*VertexLevigo) {
	return edge.object
}

func (edge *EdgeLevigo) String() (string) {
	str := fmt.Sprintf("<EdgeLevigo:%v,s=%v,o=%v,l=%v@%v>",edge.IdAsString(), 
		edge.subject.IdAsString(),edge.object.IdAsString(),edge.label,edge.db)
	return str
}
