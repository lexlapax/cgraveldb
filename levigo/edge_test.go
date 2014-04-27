package levigo

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)


type EdgeTestSuite struct {
	coretest.EdgeTestSuite
}


func (suite *EdgeTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
    suite.TestGraph.Open(testingdir) 
}

func (suite *EdgeTestSuite) TearSuite() {
    suite.TestGraph.Close()
    suite.TestGraph = nil
}

func (suite *EdgeTestSuite) SetupTest() {
    suite.TestGraph.Clear()
}

func (suite *EdgeTestSuite) TearDownTest() {
	//   suite.TestGraph.Close()
}

func TestEdgeTestSuite(t *testing.T) {
	//t.Skip()
    suite.Run(t, new(EdgeTestSuite))
}
