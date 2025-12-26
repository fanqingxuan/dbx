package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/fanqingxuan/dbx/internal/config"
	"github.com/fanqingxuan/dbx/internal/generator"
	"github.com/fanqingxuan/dbx/internal/schema"
)

func main() {
	var (
		dsn     = flag.String("dsn", "", "MySQL DSN")
		output  = flag.String("o", ".", "Output directory")
		cfgFile = flag.String("c", "dbx.yml", "Config file path")
	)
	flag.Parse()
	tables := flag.Args()

	cfg := &config.Config{OutputDir: *output}
	if *dsn != "" {
		cfg.DSN = *dsn
	} else if f, err := config.Load(*cfgFile); err == nil {
		cfg = f
		if *output != "." {
			cfg.OutputDir = *output
		}
	}

	if cfg.DSN == "" {
		fmt.Fprintln(os.Stderr, "dsn is required")
		os.Exit(1)
	}

	// 自动检测包名
	if cfg.Package == "" {
		cfg.Package = detectPackage(cfg.OutputDir)
	}
	if cfg.Package == "" {
		fmt.Fprintln(os.Stderr, "cannot detect package, please create go.mod or specify in config")
		os.Exit(1)
	}

	if len(tables) > 0 {
		cfg.Tables = tables
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	var tbls []schema.Table
	if len(cfg.Tables) > 0 {
		tbls, err = schema.LoadTables(db, cfg.Tables)
	} else {
		tbls, err = schema.LoadAllTables(db)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gen := &generator.Generator{OutputDir: cfg.OutputDir, Package: cfg.Package}
	for _, t := range tbls {
		if err := gen.GenerateModel(t); err != nil {
			fmt.Fprintf(os.Stderr, "generate model %s: %v\n", t.Name, err)
			continue
		}
		if err := gen.GenerateGenDAO(t); err != nil {
			fmt.Fprintf(os.Stderr, "generate gen dao %s: %v\n", t.Name, err)
			continue
		}
		if err := gen.GenerateDAO(t); err != nil {
			fmt.Fprintf(os.Stderr, "generate dao %s: %v\n", t.Name, err)
			continue
		}
		fmt.Printf("generated: %s\n", t.Name)
	}
}

func detectPackage(outputDir string) string {
	absOutput, _ := filepath.Abs(outputDir)
	dir := absOutput
	for {
		modFile := filepath.Join(dir, "go.mod")
		if f, err := os.Open(modFile); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "module ") {
					mod := strings.TrimSpace(strings.TrimPrefix(line, "module"))
					rel, _ := filepath.Rel(dir, absOutput)
					if rel == "." {
						return mod
					}
					return mod + "/" + filepath.ToSlash(rel)
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
