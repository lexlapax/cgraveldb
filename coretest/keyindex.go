package coretest

import (
	"github.com/stretchr/testify/assert"
	//"github.com/lexlapax/graveldb/core"
)

type GraphKeyIndexTestSuite struct {
	BaseTestSuite
}

func (suite *GraphKeyIndexTestSuite) TestGraphEmpty(){

	//todo - graph capabilites based 
	vertices, _ := suite.TestGraph.Vertices()
	assert.Equal(suite.T(), 0, len(vertices))
	edges, _ := suite.TestGraph.Edges()
	assert.Equal(suite.T(), 0, len(edges))
}
