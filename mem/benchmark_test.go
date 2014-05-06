package mem

import (
	"testing"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)


func BenchmarkAddKnown(b *testing.B){
	graph := core.GetGraph(GraphImpl)
	graph.Open()
	coretest.SuiteBenchMark(b, graph)
	graph.Close()
}

func init() {
	Register()
}
