package mem

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

type BaseVertexQueryTestSuite struct {
	coretest.BaseVertexQueryTestSuite
}

func (suite *BaseVertexQueryTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
}

func (suite *BaseVertexQueryTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestQueryGraphTestSuite(t *testing.T) {
	t.Skip()
    suite.Run(t, new(BaseVertexQueryTestSuite))
}
