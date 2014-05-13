package coretest

import (
	// "github.com/stretchr/testify/assert"
)

type BaseGraphQueryTestSuite struct {
	BaseTestSuite
}

func (suite *BaseGraphQueryTestSuite) TestBaseGraphQuery(){

	// this is what we will test
	// v1 points to v2 and v3 which both point to v4 which points back to v1

	// 		v2
	// 	/		\
	// v1			v4 - v1
	// 	\		/
	// 		v3
	
	//

	vid1 := "vertex1"
	vid2 := "vertex2"
	vid3 := "vertex3"
	vid4 := "vertex4"

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	vertex3,_ := suite.TestGraph.AddVertex(vid3)
	vertex4,_ := suite.TestGraph.AddVertex(vid4)
	vertex1.SetProperty("name", []byte("vertex 1"))
	vertex2.SetProperty("name", []byte("vertex 2"))
	vertex3.SetProperty("name", []byte("vertex 3"))
	vertex4.SetProperty("name", []byte("vertex 4"))

	eid1 := "edge1"
	eid2 := "edge2"
	eid3 := "edge3"
	eid4 := "edge4"
	eid5 := "edge5"

	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex1, vertex3, "1 to 3")
	edge3, _ := suite.TestGraph.AddEdge(eid3, vertex2, vertex4, "2 to 4")
	edge4, _ := suite.TestGraph.AddEdge(eid4, vertex3, vertex4, "3 to 4")
	edge5, _ := suite.TestGraph.AddEdge(eid5, vertex4, vertex1, "4 to 1")
	edge1.SetProperty("name", []byte("edge 1"))
	edge2.SetProperty("name", []byte("edge 2"))
	edge3.SetProperty("name", []byte("edge 3"))
	edge4.SetProperty("name", []byte("edge 4"))
	edge5.SetProperty("name", []byte("edge 5"))

}