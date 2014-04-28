package mem

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
}

func (suite *AtomTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestAtomTestSuite(t *testing.T) {
	//t.Skip()
    suite.Run(t, new(AtomTestSuite))
}
