package parser

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// FieldInfo 字段信息
type FieldInfo struct {
	Name         string // 字段名
	Type         string // 字段类型
	IsNull       bool   // 是否允许NULL
	DefaultValue string // 默认值
	Comment      string // 注释
	IsPrimary    bool   // 是否主键
	IsAutoInc    bool   // 是否自增
}

// TableInfo 表信息
type TableInfo struct {
	Name        string      // 表名
	Comment     string      // 表注释
	Fields      []FieldInfo // 字段列表
	PrimaryKeys []string    // 主键列表
}

// SQLParser SQL解析器
type SQLParser struct{}

// NewSQLParser 创建新的SQL解析器
func NewSQLParser() *SQLParser {
	return &SQLParser{}
}

// Parse 解析SQL语句
func (p *SQLParser) Parse(sql string) (*TableInfo, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("SQL解析失败: %v", err)
	}

	ddl, ok := stmt.(*sqlparser.DDL)
	if !ok {
		return nil, fmt.Errorf("不是DDL语句")
	}

	if ddl.Action != "create" {
		return nil, fmt.Errorf("不是CREATE语句")
	}

	tableInfo := &TableInfo{
		Name:        ddl.NewName.Name.String(),
		Fields:      []FieldInfo{},
		PrimaryKeys: []string{},
	}

	// 解析表选项
	if ddl.TableSpec != nil {
		tableInfo.Comment = p.extractTableComment(ddl.TableSpec.Options)
	}

	// 解析字段
	if ddl.TableSpec != nil {
		for _, col := range ddl.TableSpec.Columns {
			field := p.parseColumn(col)
			tableInfo.Fields = append(tableInfo.Fields, field)
		}

		// 解析主键
		for _, idx := range ddl.TableSpec.Indexes {
			if idx.Info.Primary {
				for _, col := range idx.Columns {
					tableInfo.PrimaryKeys = append(tableInfo.PrimaryKeys, col.Column.String())
				}
			}
		}

		// 设置主键标识
		for i := range tableInfo.Fields {
			for _, pk := range tableInfo.PrimaryKeys {
				if tableInfo.Fields[i].Name == pk {
					tableInfo.Fields[i].IsPrimary = true
					break
				}
			}
		}
	}

	return tableInfo, nil
}

// parseColumn 解析列定义
func (p *SQLParser) parseColumn(col *sqlparser.ColumnDefinition) FieldInfo {
	field := FieldInfo{
		Name:   col.Name.String(),
		Type:   strings.ToLower(col.Type.Type),
		IsNull: true, // 默认为true
	}

	// 解析类型属性
	if col.Type.NotNull {
		field.IsNull = false
	}
	if col.Type.Autoincrement {
		field.IsAutoInc = true
	}

	// 解析默认值
	if col.Type.Default != nil {
		if col.Type.Default.Type == sqlparser.StrVal {
			field.DefaultValue = string(col.Type.Default.Val)
		} else if col.Type.Default.Type == sqlparser.IntVal {
			field.DefaultValue = string(col.Type.Default.Val)
		} else if col.Type.Default.Type == sqlparser.FloatVal {
			field.DefaultValue = string(col.Type.Default.Val)
		}
	}

	// 解析注释
	if col.Type.Comment != nil {
		field.Comment = string(col.Type.Comment.Val)
	}

	return field
}

// extractTableComment 提取表注释
func (p *SQLParser) extractTableComment(options string) string {
	// 由于 TableSpec.Options 是 string 类型，我们需要解析它
	// 这里简化处理，直接查找 COMMENT 关键字
	if strings.Contains(strings.ToUpper(options), "COMMENT") {
		// 简单的注释提取逻辑
		commentIndex := strings.Index(strings.ToUpper(options), "COMMENT")
		if commentIndex != -1 {
			commentPart := options[commentIndex:]
			// 查找引号包围的注释
			if strings.Contains(commentPart, "'") {
				start := strings.Index(commentPart, "'")
				end := strings.LastIndex(commentPart, "'")
				if start != -1 && end != -1 && start < end {
					return commentPart[start+1 : end]
				}
			}
		}
	}
	return ""
}
