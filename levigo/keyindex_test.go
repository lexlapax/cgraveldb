package levigo

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

type GraphKeyIndexTestSuite struct {
	coretest.GraphKeyIndexTestSuite
}

func (suite *GraphKeyIndexTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
    suite.TestGraph.Open(testingdir) 
}

func (suite *GraphKeyIndexTestSuite) TearSuite() {
    suite.TestGraph.Close()
    suite.TestGraph = nil
}

func (suite *GraphKeyIndexTestSuite) SetupTest() {
    suite.TestGraph.Clear()
}

func (suite *GraphKeyIndexTestSuite) TearDownTest() {
	//   suite.TestGraph.Close()
}

func TestGraphKeyIndexTestSuite(t *testing.T) {
	// t.Skip()
    suite.Run(t, new(GraphKeyIndexTestSuite))
}
