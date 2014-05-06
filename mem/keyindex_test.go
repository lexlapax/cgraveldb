package mem

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

func init() {
	Register()
}

type GraphKeyIndexTestSuite struct {
	coretest.GraphKeyIndexTestSuite
}

func (suite *GraphKeyIndexTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
}

func (suite *GraphKeyIndexTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestGraphKeyIndexTestSuite(t *testing.T) {
	t.Skip()
    suite.Run(t, new(GraphKeyIndexTestSuite))
}
