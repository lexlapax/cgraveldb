package mem

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
}

func (suite *VertexTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestVertexTestSuite(t *testing.T) {
    suite.Run(t, new(VertexTestSuite))
}
