
package core

import (
)

var graphDbs = make(map[string]Graph)

func Register(graphimpl string, graph Graph) {
	if graph == nil {
		panic("core: Register graph implementation is nil")
	}
	if _, ok := graphDbs[graphimpl]; ok {
		panic("core: Register called twice for graph implementation " + graphimpl)
	}
	graphDbs[graphimpl] = graph
}

func GetGraph(graphimpl string) Graph {
	if graph, ok := graphDbs[graphimpl]; ok {
		return graph
	} else {
		panic("core: Graph implementation not found for " + graphimpl)
	}

}