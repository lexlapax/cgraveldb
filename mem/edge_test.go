package mem

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
    suite.TestGraph = core.GetGraph(graphimpl)
}

func (suite *EdgeTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestEdgeTestSuite(t *testing.T) {
    suite.Run(t, new(EdgeTestSuite))
}
