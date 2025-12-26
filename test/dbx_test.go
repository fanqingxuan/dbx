package test

import (
	"testing"

	"github.com/fanqingxuan/dbx/pkg/dbx"
	"github.com/jmoiron/sqlx"
)

type mockLogger struct {
	logs []string
}

func (m *mockLogger) Debug(format string, args ...any) { m.logs = append(m.logs, "debug") }
func (m *mockLogger) Info(format string, args ...any)  { m.logs = append(m.logs, "info") }
func (m *mockLogger) Warn(format string, args ...any)  { m.logs = append(m.logs, "warn") }
func (m *mockLogger) Error(format string, args ...any) { m.logs = append(m.logs, "error") }

func TestNew(t *testing.T) {
	sqlxDB, _ := sqlx.Open("mysql", "fake:fake@tcp(localhost:3306)/fake")
	db := dbx.New(sqlxDB)
	if db == nil {
		t.Fatal("New() returned nil")
	}
}

func TestSetLogger(t *testing.T) {
	sqlxDB, _ := sqlx.Open("mysql", "fake:fake@tcp(localhost:3306)/fake")
	db := dbx.New(sqlxDB)

	logger := &mockLogger{}
	db.SetLogger(logger)
}

func TestIn(t *testing.T) {
	sqlxDB, _ := sqlx.Open("mysql", "fake:fake@tcp(localhost:3306)/fake")
	db := dbx.New(sqlxDB)

	query, args, err := db.In("SELECT * FROM user WHERE id IN (?)", []int{1, 2, 3})
	if err != nil {
		t.Fatalf("In() error: %v", err)
	}
	if query != "SELECT * FROM user WHERE id IN (?, ?, ?)" {
		t.Fatalf("In() query = %s, want SELECT * FROM user WHERE id IN (?, ?, ?)", query)
	}
	if len(args) != 3 {
		t.Fatalf("In() args len = %d, want 3", len(args))
	}
}
