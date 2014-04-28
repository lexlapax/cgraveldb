package coretest

import (
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
)

type BaseTestSuite struct {
	suite.Suite
	TestGraph core.Graph
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *BaseTestSuite) SetupSuite() {
    suite.TestGraph = nil
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *BaseTestSuite) TearSuite() {
    suite.TestGraph = nil
}

// The SetupTest method will be run before every test in the suite.
func (suite *BaseTestSuite) SetupTest() {
    suite.TestGraph.Open()
    suite.TestGraph.Clear()
}

// The TearDownTest method will be run after every test in the suite.
func (suite *BaseTestSuite) TearDownTest() {
    suite.TestGraph.Close()
}

