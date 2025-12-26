package schema

import (
	"database/sql"
	"strings"
)

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name       string
	Type       string
	Nullable   bool
	IsPrimary  bool
	Comment    string
}

func (c Column) GoName() string {
	return toCamelCase(c.Name, true)
}

func (c Column) GoType() string {
	base := mysqlToGoType(c.Type)
	if c.Nullable && !strings.HasPrefix(base, "*") && base != "any" {
		return "*" + base
	}
	return base
}

func (t Table) GoName() string {
	return toCamelCase(t.Name, true)
}

func (t Table) PrimaryKey() *Column {
	for _, c := range t.Columns {
		if c.IsPrimary {
			return &c
		}
	}
	return nil
}

func LoadTables(db *sql.DB, tables []string) ([]Table, error) {
	var result []Table
	for _, name := range tables {
		t, err := loadTable(db, name)
		if err != nil {
			return nil, err
		}
		result = append(result, *t)
	}
	return result, nil
}

func LoadAllTables(db *sql.DB) ([]Table, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return LoadTables(db, tables)
}

func loadTable(db *sql.DB, name string) (*Table, error) {
	rows, err := db.Query("SHOW FULL COLUMNS FROM " + name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	t := &Table{Name: name}
	for rows.Next() {
		var field, colType, null, key string
		var collation, extra, privileges, comment sql.NullString
		var defaultVal sql.NullString

		if err := rows.Scan(&field, &colType, &collation, &null, &key, &defaultVal, &extra, &privileges, &comment); err != nil {
			return nil, err
		}

		t.Columns = append(t.Columns, Column{
			Name:      field,
			Type:      colType,
			Nullable:  null == "YES",
			IsPrimary: key == "PRI",
			Comment:   comment.String,
		})
	}
	return t, rows.Err()
}

func toCamelCase(s string, upper bool) string {
	var b strings.Builder
	nextUpper := upper
	for _, c := range s {
		if c == '_' {
			nextUpper = true
			continue
		}
		if nextUpper {
			b.WriteRune(toUpper(c))
			nextUpper = false
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

func mysqlToGoType(t string) string {
	t = strings.ToLower(t)
	switch {
	case strings.HasPrefix(t, "tinyint(1)"):
		return "bool"
	case strings.HasPrefix(t, "tinyint"):
		return "int8"
	case strings.HasPrefix(t, "smallint"):
		return "int16"
	case strings.HasPrefix(t, "mediumint"), strings.HasPrefix(t, "int"):
		return "int32"
	case strings.HasPrefix(t, "bigint"):
		return "int64"
	case strings.HasPrefix(t, "float"):
		return "float32"
	case strings.HasPrefix(t, "double"), strings.HasPrefix(t, "decimal"):
		return "float64"
	case strings.HasPrefix(t, "varchar"), strings.HasPrefix(t, "char"),
		strings.HasPrefix(t, "text"), strings.HasPrefix(t, "mediumtext"),
		strings.HasPrefix(t, "longtext"), strings.HasPrefix(t, "enum"):
		return "string"
	case strings.HasPrefix(t, "datetime"), strings.HasPrefix(t, "timestamp"):
		return "time.Time"
	case strings.HasPrefix(t, "date"):
		return "time.Time"
	case strings.HasPrefix(t, "blob"), strings.HasPrefix(t, "binary"), strings.HasPrefix(t, "varbinary"):
		return "[]byte"
	case strings.HasPrefix(t, "json"):
		return "any"
	default:
		return "string"
	}
}
