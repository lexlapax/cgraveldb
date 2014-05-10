package coretest

import (
	"fmt"
	"testing"
	"github.com/lexlapax/graveldb/core"
)

func SuiteBenchMark(b *testing.B, graph core.Graph) {
	BenchmarkAddKnown(b, graph)
}

func BenchmarkAddKnown(b *testing.B, graph core.Graph){
	graph.Clear()
	vertexprefix := "vertex-"
	edgeprefix := "edge-"
	labelprefix := "label-"
	var fromvertex, tovertex core.Vertex
	var fromvertexint, tovertexint int
	var err error
	//upto := b.N
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		upto := b.N
		fromvertexint = n + 1 
		tovertexint = n + 2
		fromvertex, err = graph.AddVertex(fmt.Sprintf("%v%v", vertexprefix, fromvertexint))
		if n == upto - 1  {
			tovertexint = 1
			tovertex, _ = graph.Vertex(fmt.Sprintf("%v%v", vertexprefix, tovertexint))
		} else {
			tovertex, _ = graph.AddVertex(fmt.Sprintf("%v%v", vertexprefix, tovertexint))
		}
		fmt.Printf("f=%v,t=%v,e=%v\n", fromvertex, tovertex,err)
		graph.AddEdge(fmt.Sprintf("%v%v", edgeprefix, n), fromvertex, tovertex, fmt.Sprintf("%v%v to %v", labelprefix, fromvertexint, tovertexint ))
	}
	fmt.Printf("vertices=%v, edges=%v\n", graph.VertexCount(), graph.EdgeCount())
}


// func BenchmarkAddKnown(b *testing.B, graph core.Graph){
// 	vertexprefix := "vertex-"
// 	edgeprefix := "edge-"
// 	//upto := b.N
// 	b.ResetTimer()
// 	for n := 1; n <= b.N; n++ {
// 		fromvertex, _ := graph.AddVertex([]byte(fmt.Sprintf("%v%v", vertexprefix, i)))
// 		tovertexint := 0
// 		if n == b.N {
// 			tovertexint =  1
// 		} else {
// 			tovertexint =  n + 1
// 		}
// 		tovertex,_ := graph.Vertex([]byte(fmt.Sprintf("%v%v", vertexprefix, tovertexint)))
// 		graph.AddEdge([]byte(fmt.Sprintf("%v%v", edgeprefix, i)), fromvertex, tovertex, fmt.Sprintf("label %v to %v", i, tovertexint ))
// 	}
// }

