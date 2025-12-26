package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanqingxuan/dbx/internal/generator"
	"github.com/fanqingxuan/dbx/internal/schema"
)

func TestGenerateModel(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dbx-test-*")
	defer os.RemoveAll(tmpDir)

	g := &generator.Generator{OutputDir: tmpDir, Package: "github.com/test/project"}

	tbl := schema.Table{
		Name: "user",
		Columns: []schema.Column{
			{Name: "id", Type: "int(11)", IsPrimary: true},
			{Name: "name", Type: "varchar(255)", Nullable: false},
			{Name: "age", Type: "int(11)", Nullable: true},
		},
	}

	if err := g.GenerateModel(tbl); err != nil {
		t.Fatalf("GenerateModel() error: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(tmpDir, "model", "user.go"))
	s := string(content)

	checks := []string{"type User struct", "Id", "Name", "*int32", "TableName()"}
	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Errorf("missing: %s", c)
		}
	}
}

func TestGenerateGenDAO(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dbx-test-*")
	defer os.RemoveAll(tmpDir)

	g := &generator.Generator{OutputDir: tmpDir, Package: "github.com/test/project"}

	tbl := schema.Table{
		Name: "user",
		Columns: []schema.Column{
			{Name: "id", Type: "int(11)", IsPrimary: true},
			{Name: "name", Type: "varchar(255)"},
		},
	}

	if err := g.GenerateGenDAO(tbl); err != nil {
		t.Fatalf("GenerateGenDAO() error: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(tmpDir, "dao", "gen", "user.go"))
	s := string(content)

	checks := []string{
		"type UserGen struct",
		"func NewUserGen",
		"Insert", "Update", "Delete",
		"FindByID", "FindByIds",
		"DeleteByIds", "UpdateByIds",
		"QueryRowsCtx", "QueryRowCtx", "ExecCtx",
	}
	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Errorf("missing: %s", c)
		}
	}
}

func TestGenerateDAO(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dbx-test-*")
	defer os.RemoveAll(tmpDir)

	g := &generator.Generator{OutputDir: tmpDir, Package: "github.com/test/project"}
	tbl := schema.Table{Name: "user"}

	g.GenerateDAO(tbl)

	content, _ := os.ReadFile(filepath.Join(tmpDir, "dao", "user.go"))
	s := string(content)

	if !strings.Contains(s, "type UserDAO struct") {
		t.Error("missing struct")
	}
	if !strings.Contains(s, "*gen.UserGen") {
		t.Error("missing embedded gen")
	}
}

func TestGenerateDAONotOverwrite(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dbx-test-*")
	defer os.RemoveAll(tmpDir)

	g := &generator.Generator{OutputDir: tmpDir, Package: "github.com/test/project"}
	tbl := schema.Table{Name: "user"}

	g.GenerateDAO(tbl)

	daoPath := filepath.Join(tmpDir, "dao", "user.go")
	os.WriteFile(daoPath, []byte("custom"), 0644)

	g.GenerateDAO(tbl)

	content, _ := os.ReadFile(daoPath)
	if string(content) != "custom" {
		t.Error("should not overwrite")
	}
}
