package main_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/stretchr/testify/suite"

	"github.com/chiselstrike/libsql-shell/src/lib"
	"github.com/chiselstrike/libsql-shell/testing/utils"
)

type RootCommandShellSuite struct {
	suite.Suite

	// dbInTest utils.DbType
	dbPath string
	tc     *utils.DbTestContext
}

func NewRootCommandShellSuite(dbPath string) *RootCommandShellSuite {
	return &RootCommandShellSuite{dbPath: dbPath}
}

func (s *RootCommandShellSuite) SetupTest() {
	s.tc = utils.NewTestContext(s.T(), s.dbPath)
}

func (s *RootCommandShellSuite) TearDownTest() {
	s.tc.TearDown()
}

func (s *RootCommandShellSuite) Test_WhenCreateTable_ExpectDbHaveTheTable() {
	outS, errS, err := s.tc.ExecuteShell([]string{"CREATE TABLE test (name STRING);", "SELECT * FROM test;"})
	s.tc.Assert(err, qt.IsNil)
	s.tc.Assert(errS, qt.Equals, "")

	s.tc.Assert(outS, qt.Equals, utils.GetPrintTableOutput([]string{"name"}, [][]string{}))
}

func (s *RootCommandShellSuite) Test_WhenCreateTableAndInsertData_ExpectDbHaveTheTableWithTheData() {
	outS, errS, err := s.tc.ExecuteShell([]string{"CREATE TABLE test (name STRING);", "INSERT INTO test VALUES ('test');", "SELECT * FROM test;"})
	s.tc.Assert(err, qt.IsNil)
	s.tc.Assert(errS, qt.Equals, "")

	s.tc.Assert(outS, qt.Equals, utils.GetPrintTableOutput([]string{"name"}, [][]string{{"test"}}))
}

func (s *RootCommandShellSuite) Test_WhenNoCommandsAreProvided_ExpectShellExecutedWithoutError() {
	outS, errS, err := s.tc.ExecuteShell([]string{})

	s.tc.Assert(err, qt.IsNil)
	s.tc.Assert(errS, qt.Equals, "")
	s.tc.Assert(outS, qt.Equals, "")
}

func (s *RootCommandShellSuite) Test_WhenExecuteInvalidStatement_ExpectError() {
	outS, errS, err := s.tc.ExecuteShell([]string{"SELECTT 1;"})
	s.tc.Assert(err, qt.IsNil)

	s.tc.Assert(outS, qt.Equals, "")
	s.tc.Assert(len(errS), qt.Not(qt.Equals), 0)
}

func (s *RootCommandShellSuite) Test_WhenTypingQuitCommand_ExpectShellNotRunFollowingCommands() {
	outS, errS, err := s.tc.ExecuteShell([]string{lib.QUIT_COMMAND, "SELECT 1;"})

	s.tc.Assert(err, qt.IsNil)
	s.tc.Assert(outS, qt.Equals, "")
	s.tc.Assert(errS, qt.Equals, "")
}

func TestRootCommandShellSuite_WhenDbIsSQLite(t *testing.T) {
	suite.Run(t, NewRootCommandShellSuite(t.TempDir()+"test.sqlite"))
}

func TestRootCommandShellSuite_WhenDbIsTurso(t *testing.T) {
	testConfig := utils.GetTestConfig(t)
	if testConfig.SkipTursoTests {
		t.Skip("Skipping Turso tests due configuration")
	}

	suite.Run(t, NewRootCommandShellSuite(testConfig.TursoDbPath))
}
