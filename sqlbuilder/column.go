// Modeling of columns

package sqlbuilder

import (
	"strings"
)

type Column interface {
	Expression

	Name() string
	TableName() string

	DefaultAlias() Projection
	// Internal function for tracking tableName that a column belongs to
	// for the purpose of serialization
	setTableName(table string)
}

type NullableColumn bool

const (
	Nullable    NullableColumn = true
	NotNullable NullableColumn = false
)

type Collation string

const (
	UTF8CaseInsensitive Collation = "utf8_unicode_ci"
	UTF8CaseSensitive   Collation = "utf8_unicode"
	UTF8Binary          Collation = "utf8_bin"
)

// Representation of MySQL charsets
type Charset string

const (
	UTF8 Charset = "utf8"
)

// The base type for real materialized columns.
type baseColumn struct {
	expressionInterfaceImpl

	name      string
	nullable  NullableColumn
	tableName string
}

func newBaseColumn(name string, nullable NullableColumn, tableName string, parent Column) baseColumn {
	bc := baseColumn{
		name:      name,
		nullable:  nullable,
		tableName: tableName,
	}

	bc.expressionInterfaceImpl.parent = parent

	return bc
}

func (c *baseColumn) Name() string {
	return c.name
}

func (c *baseColumn) TableName() string {
	return c.tableName
}

func (c *baseColumn) setTableName(table string) {
	c.tableName = table
}

func (c *baseColumn) DefaultAlias() Projection {
	return c.AS(c.tableName + "." + c.name)
}

func (c baseColumn) Serialize(out *queryData, options ...serializeOption) error {

	setOrderBy := out.statementType == set_statement && out.clauseType == order_by_clause

	if setOrderBy {
		out.WriteString(`"`)
	}

	if c.tableName != "" {
		out.WriteString(c.tableName)
		out.WriteString(".")
	}

	wrapColumnName := strings.Contains(c.name, ".") && !setOrderBy

	if wrapColumnName {
		out.WriteString(`"`)
	}

	out.WriteString(c.name)

	if wrapColumnName {
		out.WriteString(`"`)
	}

	if setOrderBy {
		out.WriteString(`"`)
	}

	return nil
}
