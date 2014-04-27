package levigo

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

func init() {
	Register()
}

var testingdir = "./testing.db"
type GraphTestSuite struct {
	coretest.GraphTestSuite
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *GraphTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
    suite.TestGraph.Open(testingdir) 
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *GraphTestSuite) TearSuite() {
    suite.TestGraph.Close()
    suite.TestGraph = nil
}

// The SetupTest method will be run before every test in the suite.
func (suite *GraphTestSuite) SetupTest() {
    suite.TestGraph.Clear()
}

// The TearDownTest method will be run after every test in the suite.
func (suite *GraphTestSuite) TearDownTest() {
	//   suite.TestGraph.Close()
}


func TestGraphTestSuite(t *testing.T) {
	//t.Skip()
    suite.Run(t, new(GraphTestSuite))
}
