package database

import (
	"github.com/codebee2/sql2gorm/parser"
)

// ToParserTableInfo 将数据库表结构信息转换为解析器的表信息格式
func (ti *TableInfo) ToParserTableInfo() *parser.TableInfo {
	fields := make([]parser.FieldInfo, len(ti.Columns))

	for i, col := range ti.Columns {
		fields[i] = parser.FieldInfo{
			Name:         col.Name,
			Type:         col.Type,
			IsNull:       col.IsNullable,
			DefaultValue: col.DefaultValue,
			Comment:      col.Comment,
			IsPrimary:    col.IsPrimaryKey,
			IsAutoInc:    col.IsAutoIncr,
		}
	}

	return &parser.TableInfo{
		Name:        ti.Name,
		Comment:     ti.Comment,
		Fields:      fields,
		PrimaryKeys: ti.PrimaryKeys,
	}
}
