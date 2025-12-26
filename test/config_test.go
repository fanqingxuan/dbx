package test

import (
	"os"
	"testing"

	"github.com/fanqingxuan/dbx/internal/config"
)

func TestConfigLoad(t *testing.T) {
	content := `dsn: user:pass@tcp(localhost:3306)/test
output: ./out
package: github.com/test/project
tables: user, order, product`

	tmpFile, err := os.CreateTemp("", "dbx-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Close()

	cfg, err := config.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.DSN != "user:pass@tcp(localhost:3306)/test" {
		t.Errorf("DSN = %q, want user:pass@tcp(localhost:3306)/test", cfg.DSN)
	}
	if cfg.OutputDir != "./out" {
		t.Errorf("OutputDir = %q, want ./out", cfg.OutputDir)
	}
	if cfg.Package != "github.com/test/project" {
		t.Errorf("Package = %q, want github.com/test/project", cfg.Package)
	}
	if len(cfg.Tables) != 3 {
		t.Errorf("Tables len = %d, want 3", len(cfg.Tables))
	}
}

func TestConfigDefaults(t *testing.T) {
	content := `dsn: test`

	tmpFile, err := os.CreateTemp("", "dbx-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Close()

	cfg, _ := config.Load(tmpFile.Name())

	if cfg.OutputDir != "." {
		t.Errorf("OutputDir default = %q, want .", cfg.OutputDir)
	}
	if cfg.Package != "model" {
		t.Errorf("Package default = %q, want model", cfg.Package)
	}
}

func TestConfigNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/file.yml")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}
}
