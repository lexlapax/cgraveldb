package levigo

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

type VertexTestSuite struct {
	coretest.VertexTestSuite
}

func (suite *VertexTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
    suite.TestGraph.Open(testingdir) 
}

func (suite *VertexTestSuite) TearSuite() {
    suite.TestGraph.Close()
    suite.TestGraph = nil
}

func (suite *VertexTestSuite) SetupTest() {
    suite.TestGraph.Clear()
}

func (suite *VertexTestSuite) TearDownTest() {
	//   suite.TestGraph.Close()
}

func TestVertexTestSuite(t *testing.T) {
	t.Skip()
    suite.Run(t, new(VertexTestSuite))
}
