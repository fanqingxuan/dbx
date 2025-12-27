package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanqingxuan/dbx/internal/schema"
)

type Generator struct {
	OutputDir string
	Package   string
}

func (g *Generator) GenerateModel(t schema.Table) error {
	dir := filepath.Join(g.OutputDir, "model")
	os.MkdirAll(dir, 0755)

	var b strings.Builder
	b.WriteString("package model\n\n")

	needTime := false
	for _, c := range t.Columns {
		if strings.Contains(c.GoType(), "time.Time") {
			needTime = true
			break
		}
	}
	if needTime {
		b.WriteString("import \"time\"\n\n")
	}

	// 计算对齐宽度
	maxName, maxType := 0, 0
	for _, c := range t.Columns {
		if l := len(c.GoName()); l > maxName {
			maxName = l
		}
		if l := len(c.GoType()); l > maxType {
			maxType = l
		}
	}

	b.WriteString(fmt.Sprintf("type %s struct {\n", t.GoName()))
	for _, c := range t.Columns {
		b.WriteString(fmt.Sprintf("\t%-*s %-*s `db:\"%s\"`\n", maxName, c.GoName(), maxType, c.GoType(), c.Name))
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("func (%s) TableName() string { return \"%s\" }\n", t.GoName(), t.Name))

	return os.WriteFile(filepath.Join(dir, t.Name+".go"), []byte(b.String()), 0644)
}

func (g *Generator) GenerateDO(t schema.Table) error {
	dir := filepath.Join(g.OutputDir, "do")
	os.MkdirAll(dir, 0755)

	var b strings.Builder
	b.WriteString("package do\n\n")

	b.WriteString(fmt.Sprintf("type %s struct {\n", t.GoName()))
	for _, c := range t.Columns {
		b.WriteString(fmt.Sprintf("\t%s any\n", c.GoName()))
	}
	b.WriteString("}\n")

	return os.WriteFile(filepath.Join(dir, t.Name+".go"), []byte(b.String()), 0644)
}

func (g *Generator) GenerateGenDAO(t schema.Table) error {
	dir := filepath.Join(g.OutputDir, "dao", "gen")
	os.MkdirAll(dir, 0755)

	pk := t.PrimaryKey()
	pkType := "int64"
	pkName := "id"
	if pk != nil {
		pkType = pk.GoType()
		pkName = pk.Name
	}

	var b strings.Builder
	b.WriteString("package gen\n\n")
	b.WriteString("import (\n\t\"context\"\n\t\"strings\"\n\n\t\"github.com/fanqingxuan/dbx/pkg/dbx\"\n")
	b.WriteString(fmt.Sprintf("\t\"%s/do\"\n", g.Package))
	b.WriteString(fmt.Sprintf("\t\"%s/model\"\n)\n\n", g.Package))

	structName := t.GoName() + "Gen"
	b.WriteString(fmt.Sprintf("type %s struct {\n\tdb *dbx.DB\n}\n\n", structName))
	b.WriteString(fmt.Sprintf("func New%s(db *dbx.DB) *%s {\n\treturn &%s{db: db}\n}\n\n", structName, structName, structName))

	// Insert (使用 DO，忽略 nil 字段)
	b.WriteString(fmt.Sprintf("func (d *%s) Insert(ctx context.Context, m do.%s) error {\n", structName, t.GoName()))
	b.WriteString("\tcols, vals, args := d.nonNilFields(m)\n")
	b.WriteString("\tif len(cols) == 0 { return nil }\n")
	b.WriteString(fmt.Sprintf("\tquery := \"INSERT INTO %s (\" + strings.Join(cols, \",\") + \") VALUES (\" + strings.Join(vals, \",\") + \")\"\n", t.Name))
	b.WriteString("\t_, err := d.db.ExecCtx(ctx, query, args...)\n\treturn err\n}\n\n")

	// Delete
	b.WriteString(fmt.Sprintf("func (d *%s) Delete(ctx context.Context, %s %s) error {\n", structName, pkName, pkType))
	b.WriteString(fmt.Sprintf("\t_, err := d.db.ExecCtx(ctx, \"DELETE FROM %s WHERE %s=?\", %s)\n\treturn err\n}\n\n", t.Name, pkName, pkName))

	// FindByID
	b.WriteString(fmt.Sprintf("func (d *%s) FindByID(ctx context.Context, %s %s) (*model.%s, error) {\n", structName, pkName, pkType, t.GoName()))
	b.WriteString(fmt.Sprintf("\tvar m model.%s\n", t.GoName()))
	b.WriteString(fmt.Sprintf("\terr := d.db.QueryRowCtx(ctx, &m, \"SELECT * FROM %s WHERE %s=?\", %s)\n", t.Name, pkName, pkName))
	b.WriteString("\tif err != nil { return nil, err }\n\treturn &m, nil\n}\n\n")

	// FindByIds
	b.WriteString(fmt.Sprintf("func (d *%s) FindByIds(ctx context.Context, ids []%s) ([]*model.%s, error) {\n", structName, pkType, t.GoName()))
	b.WriteString(fmt.Sprintf("\tvar list []*model.%s\n", t.GoName()))
	b.WriteString("\tif len(ids) == 0 { return list, nil }\n")
	b.WriteString(fmt.Sprintf("\tquery, args, _ := d.db.In(\"SELECT * FROM %s WHERE %s IN (?)\", ids)\n", t.Name, pkName))
	b.WriteString("\terr := d.db.QueryRowsCtx(ctx, &list, query, args...)\n")
	b.WriteString("\treturn list, err\n}\n\n")

	// DeleteByIds
	b.WriteString(fmt.Sprintf("func (d *%s) DeleteByIds(ctx context.Context, ids []%s) (int64, error) {\n", structName, pkType))
	b.WriteString("\tif len(ids) == 0 { return 0, nil }\n")
	b.WriteString(fmt.Sprintf("\tquery, args, _ := d.db.In(\"DELETE FROM %s WHERE %s IN (?)\", ids)\n", t.Name, pkName))
	b.WriteString("\treturn d.ExecCtx(ctx, query, args...)\n}\n\n")

	// UpdateById (使用 DO，忽略 nil 字段)
	b.WriteString(fmt.Sprintf("func (d *%s) UpdateById(ctx context.Context, m do.%s, %s %s) (int64, error) {\n", structName, t.GoName(), pkName, pkType))
	b.WriteString("\tcols, _, args := d.nonNilFields(m)\n")
	b.WriteString("\tif len(cols) == 0 { return 0, nil }\n")
	b.WriteString("\tvar sets []string\n")
	b.WriteString("\tfor _, c := range cols { sets = append(sets, c+\"=?\") }\n")
	b.WriteString(fmt.Sprintf("\targs = append(args, %s)\n", pkName))
	b.WriteString(fmt.Sprintf("\tquery := \"UPDATE %s SET \" + strings.Join(sets, \",\") + \" WHERE %s=?\"\n", t.Name, pkName))
	b.WriteString("\treturn d.ExecCtx(ctx, query, args...)\n}\n\n")

	// UpdateByIds (使用 DO，忽略 nil 字段)
	b.WriteString(fmt.Sprintf("func (d *%s) UpdateByIds(ctx context.Context, m do.%s, ids []%s) (int64, error) {\n", structName, t.GoName(), pkType))
	b.WriteString("\tif len(ids) == 0 { return 0, nil }\n")
	b.WriteString("\tcols, _, args := d.nonNilFields(m)\n")
	b.WriteString("\tif len(cols) == 0 { return 0, nil }\n")
	b.WriteString("\tvar sets []string\n")
	b.WriteString("\tfor _, c := range cols { sets = append(sets, c+\"=?\") }\n")
	b.WriteString(fmt.Sprintf("\tquery := \"UPDATE %s SET \" + strings.Join(sets, \",\") + \" WHERE %s IN (?)\"\n", t.Name, pkName))
	b.WriteString("\tquery, inArgs, _ := d.db.In(query, ids)\n")
	b.WriteString("\targs = append(args, inArgs...)\n")
	b.WriteString("\treturn d.ExecCtx(ctx, query, args...)\n}\n\n")

	// 暴露查询方法
	b.WriteString(fmt.Sprintf("func (d *%s) QueryRowsCtx(ctx context.Context, dest any, query string, args ...any) error {\n", structName))
	b.WriteString("\treturn d.db.QueryRowsCtx(ctx, dest, query, args...)\n}\n\n")
	b.WriteString(fmt.Sprintf("func (d *%s) QueryRowCtx(ctx context.Context, dest any, query string, args ...any) error {\n", structName))
	b.WriteString("\treturn d.db.QueryRowCtx(ctx, dest, query, args...)\n}\n\n")
	b.WriteString(fmt.Sprintf("func (d *%s) QueryValueCtx(ctx context.Context, dest any, query string, args ...any) error {\n", structName))
	b.WriteString("\treturn d.db.QueryValueCtx(ctx, dest, query, args...)\n}\n\n")
	b.WriteString(fmt.Sprintf("func (d *%s) QueryColumnCtx(ctx context.Context, dest any, query string, args ...any) error {\n", structName))
	b.WriteString("\treturn d.db.QueryColumnCtx(ctx, dest, query, args...)\n}\n\n")
	b.WriteString(fmt.Sprintf("func (d *%s) ExecCtx(ctx context.Context, query string, args ...any) (int64, error) {\n", structName))
	b.WriteString("\tresult, err := d.db.ExecCtx(ctx, query, args...)\n\tif err != nil { return 0, err }\n\treturn result.RowsAffected()\n}\n\n")

	// nonNilFields helper (遍历 DO 结构体)
	b.WriteString(fmt.Sprintf("func (d *%s) nonNilFields(m do.%s) ([]string, []string, []any) {\n", structName, t.GoName()))
	b.WriteString("\tvar cols, vals []string\n\tvar args []any\n")
	for _, c := range t.Columns {
		b.WriteString(fmt.Sprintf("\tif m.%s != nil { cols = append(cols, \"%s\"); vals = append(vals, \"?\"); args = append(args, m.%s) }\n", c.GoName(), c.Name, c.GoName()))
	}
	b.WriteString("\treturn cols, vals, args\n}\n")

	return os.WriteFile(filepath.Join(dir, t.Name+".go"), []byte(b.String()), 0644)
}

func (g *Generator) GenerateDAO(t schema.Table) error {
	dir := filepath.Join(g.OutputDir, "dao")
	os.MkdirAll(dir, 0755)

	path := filepath.Join(dir, t.Name+".go")
	if _, err := os.Stat(path); err == nil {
		return nil // 文件已存在，不覆盖
	}

	var b strings.Builder
	b.WriteString("package dao\n\n")
	b.WriteString("import (\n\t\"github.com/fanqingxuan/dbx/pkg/dbx\"\n")
	b.WriteString(fmt.Sprintf("\t\"%s/dao/gen\"\n)\n\n", g.Package))

	structName := t.GoName() + "DAO"
	b.WriteString(fmt.Sprintf("type %s struct {\n\t*gen.%sGen\n}\n\n", structName, t.GoName()))
	b.WriteString(fmt.Sprintf("func New%s(db *dbx.DB) *%s {\n", structName, structName))
	b.WriteString(fmt.Sprintf("\treturn &%s{%sGen: gen.New%sGen(db)}\n}\n", structName, t.GoName(), t.GoName()))

	return os.WriteFile(path, []byte(b.String()), 0644)
}
