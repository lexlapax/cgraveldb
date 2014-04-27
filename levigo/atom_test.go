package levigo

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

type AtomTestSuite struct {
	coretest.AtomTestSuite
}


func (suite *AtomTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
    suite.TestGraph.Open(testingdir) 
}

func (suite *AtomTestSuite) TearSuite() {
    suite.TestGraph.Close()
    suite.TestGraph = nil
}

func (suite *AtomTestSuite) SetupTest() {
    suite.TestGraph.Clear()
}

func (suite *AtomTestSuite) TearDownTest() {
	//   suite.TestGraph.Close()
}

func TestAtomTestSuite(t *testing.T) {
	//t.Skip()
    suite.Run(t, new(AtomTestSuite))
}
