package levigo


import (
		"fmt"
		"github.com/lexlapax/graveldb/core"
)

type EdgeLevigo struct {
	*AtomLevigo
	subject *VertexLevigo
	object *VertexLevigo
	label string
}


func (edge *EdgeLevigo) Label() string {
	return edge.label
}

func (edge *EdgeLevigo) Vertex(direction core.Direction) (core.Vertex, error) {
	if direction == core.DirOut {
		return edge.subject, nil
	} else if direction == core.DirIn {
		return edge.object, nil
	} else {
		return nil, core.ErrDirAnyUnsupported
	}
}


func (edge *EdgeLevigo) VertexOut() (core.Vertex, error) {
	return edge.subject, nil
}

func (edge *EdgeLevigo) VertexIn() (core.Vertex, error) {
	return edge.object, nil
}

func (edge *EdgeLevigo) String() (string) {
	str := fmt.Sprintf("<EdgeLevigo:%v,s=%v,o=%v,l=%v@%v>",edge.Id(), 
		edge.subject.Id(),edge.object.Id(),edge.label,edge.db)
	return str
}

