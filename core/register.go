
package core

import (
)

var graphDbs = make(map[string]Graph)

func Register(name string, graph Graph) {
	if graph == nil {
		panic("core: Register graph implementation is nil")
	}
	if _, ok := graphDbs[name]; ok {
		panic("core: Register called twice for graph implementation " + name)
	}
	graphDbs[name] = graph
}