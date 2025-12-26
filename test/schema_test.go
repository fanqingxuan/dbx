package test

import (
	"testing"

	"github.com/fanqingxuan/dbx/internal/schema"
)

func TestColumnGoName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"user_name", "UserName"},
		{"id", "Id"},
		{"created_at", "CreatedAt"},
	}

	for _, tt := range tests {
		c := schema.Column{Name: tt.name}
		if got := c.GoName(); got != tt.want {
			t.Errorf("GoName(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestColumnGoType(t *testing.T) {
	tests := []struct {
		typ      string
		nullable bool
		want     string
	}{
		{"int(11)", false, "int32"},
		{"int(11)", true, "*int32"},
		{"bigint(20)", false, "int64"},
		{"tinyint(1)", false, "bool"},
		{"varchar(255)", false, "string"},
		{"varchar(255)", true, "*string"},
		{"datetime", false, "time.Time"},
		{"json", true, "any"},
	}

	for _, tt := range tests {
		c := schema.Column{Type: tt.typ, Nullable: tt.nullable}
		if got := c.GoType(); got != tt.want {
			t.Errorf("GoType(%q, nullable=%v) = %q, want %q", tt.typ, tt.nullable, got, tt.want)
		}
	}
}

func TestTableGoName(t *testing.T) {
	tbl := schema.Table{Name: "user_detail"}
	if got := tbl.GoName(); got != "UserDetail" {
		t.Errorf("GoName() = %q, want UserDetail", got)
	}
}

func TestTablePrimaryKey(t *testing.T) {
	tbl := schema.Table{
		Columns: []schema.Column{
			{Name: "id", IsPrimary: true},
			{Name: "name"},
		},
	}

	pk := tbl.PrimaryKey()
	if pk == nil || pk.Name != "id" {
		t.Error("PrimaryKey() failed")
	}

	tbl2 := schema.Table{Columns: []schema.Column{{Name: "name"}}}
	if tbl2.PrimaryKey() != nil {
		t.Error("should return nil")
	}
}
