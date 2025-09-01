package generator

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/codebee2/sql2gorm/config"
	"github.com/codebee2/sql2gorm/parser"
	"github.com/jinzhu/inflection"
)

// GORMGenerator GORM代码生成器
type GORMGenerator struct {
	config *config.Config
}

// NewGORMGenerator 创建新的GORM生成器
func NewGORMGenerator(cfg *config.Config) *GORMGenerator {
	return &GORMGenerator{
		config: cfg,
	}
}

// Generate 生成GORM模型代码
func (g *GORMGenerator) Generate(tableInfo *parser.TableInfo) (string, error) {
	// 生成结构体名称
	structName := g.generateStructName(tableInfo.Name)

	// 生成实例变量名
	instanceName := g.generateInstanceName(structName)

	// 生成字段代码
	fields, err := g.generateFields(tableInfo.Fields)
	if err != nil {
		return "", err
	}

	// 生成TableName方法
	tableNameMethod := g.generateTableNameMethod(tableInfo.Name)

	// 组装完整代码
	code := fmt.Sprintf(`package %s

var %s = new(%s)

// %s %s
type %s struct {
%s
}

%s`,
		g.config.PackageName,
		instanceName,
		structName,
		structName,
		tableInfo.Comment,
		structName,
		fields,
		tableNameMethod)

	// 格式化Go代码
	formattedCode, err := g.formatGoCode(code)
	if err != nil {
		return "", fmt.Errorf("代码格式化失败: %v", err)
	}

	return formattedCode, nil
}

// generateStructName 生成结构体名称
func (g *GORMGenerator) generateStructName(tableName string) string {
	// 移除表名前缀（如ny_）
	parts := strings.Split(tableName, "_")
	var nameParts []string

	for _, part := range parts[1:] {
		if part != "" {
			nameParts = append(nameParts, strings.Title(part))
		}
	}

	// 如果没有下划线，直接首字母大写
	if len(nameParts) == 0 {
		return strings.Title(tableName)
	}

	// 使用inflection库处理单复数
	baseName := strings.Join(nameParts, "")
	return inflection.Singular(baseName) + "Model"
}

// generateInstanceName 生成实例变量名
func (g *GORMGenerator) generateInstanceName(structName string) string {
	return structName + "Ins"
}

// generateFields 生成字段代码
func (g *GORMGenerator) generateFields(fields []parser.FieldInfo) (string, error) {
	var fieldLines []string

	for _, field := range fields {
		fieldLine, err := g.generateField(field)
		if err != nil {
			return "", err
		}
		fieldLines = append(fieldLines, fieldLine)
	}

	return strings.Join(fieldLines, "\n\t"), nil
}

// generateField 生成单个字段代码
func (g *GORMGenerator) generateField(field parser.FieldInfo) (string, error) {
	// 生成Go字段名
	goFieldName := g.generateGoFieldName(field.Name)

	// 生成Go类型
	goType := g.mapMySQLTypeToGoType(field.Type)

	// 生成JSON标签
	jsonTag := fmt.Sprintf(`json:"%s"`, field.Name)

	// 生成GORM标签
	gormTags := g.generateGORMTags(field)

	// 生成注释
	comment := ""
	if field.Comment != "" {
		comment = fmt.Sprintf(" // %s", field.Comment)
	}

	// 组装字段代码
	fieldCode := fmt.Sprintf("\t%s\t%s `%s %s`%s",
		goFieldName,
		goType,
		jsonTag,
		gormTags,
		comment)

	return fieldCode, nil
}

// generateGoFieldName 生成Go字段名
func (g *GORMGenerator) generateGoFieldName(dbFieldName string) string {
	parts := strings.Split(dbFieldName, "_")
	var nameParts []string

	for _, part := range parts {
		if part != "" {
			nameParts = append(nameParts, strings.Title(part))
		}
	}

	return strings.Join(nameParts, "")
}

// mapMySQLTypeToGoType 映射MySQL类型到Go类型
func (g *GORMGenerator) mapMySQLTypeToGoType(mysqlType string) string {
	// 移除长度信息
	baseType := strings.Split(mysqlType, "(")[0]

	switch baseType {
	case "int", "tinyint", "smallint", "mediumint", "bigint":
		return "int64"
	case "float", "double", "decimal":
		return "float64"
	case "varchar", "char", "text", "longtext", "mediumtext", "tinytext":
		return "string"
	case "date", "datetime", "timestamp", "time":
		return "string"
	case "blob", "longblob", "mediumblob", "tinyblob":
		return "[]byte"
	case "json":
		return "string"
	default:
		return "string"
	}
}

// generateGORMTags 生成GORM标签
func (g *GORMGenerator) generateGORMTags(field parser.FieldInfo) string {
	var tags []string

	// column标签
	tags = append(tags, fmt.Sprintf("column:%s", field.Name))

	// primaryKey标签
	if field.IsPrimary {
		tags = append(tags, "primaryKey")
	}

	// autoIncrement标签
	if field.IsAutoInc {
		tags = append(tags, "autoIncrement")
	}

	// default标签
	if field.DefaultValue != "" {
		// 处理特殊默认值
		defaultValue := field.DefaultValue
		if defaultValue == "CURRENT_TIMESTAMP" {
			tags = append(tags, "default:CURRENT_TIMESTAMP")
		} else if defaultValue == "NULL" {
			tags = append(tags, "default:null")
		} else {
			tags = append(tags, fmt.Sprintf("default:%s", defaultValue))
		}
	}

	// not null标签
	if !field.IsNull {
		tags = append(tags, "not null")
	}

	return fmt.Sprintf("gorm:\"%s\"", strings.Join(tags, ";"))
}

// generateTableNameMethod 生成TableName方法
func (g *GORMGenerator) generateTableNameMethod(tableName string) string {
	return fmt.Sprintf(`// TableName 表名称
func (own *%s) TableName() string {
	return "%s"
}`, g.generateStructName(tableName), tableName)
}

// formatGoCode 格式化Go代码
func (g *GORMGenerator) formatGoCode(code string) (string, error) {
	// 使用go/format包格式化代码
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}
