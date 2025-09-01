package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// TableInfo 表结构信息
type TableInfo struct {
	Name        string
	Comment     string
	Columns     []ColumnInfo
	PrimaryKeys []string
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name         string
	Type         string
	IsNullable   bool
	DefaultValue string
	Comment      string
	IsPrimaryKey bool
	IsAutoIncr   bool
}

// Database 数据库连接管理器
type Database struct {
	db *sql.DB
}

// NewDatabase 创建新的数据库连接
func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	return &Database{db: db}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	return d.db.Close()
}

// GetTableInfo 获取表结构信息
func (d *Database) GetTableInfo(tableName string) (*TableInfo, error) {
	// 获取表注释
	tableComment, err := d.getTableComment(tableName)
	if err != nil {
		return nil, err
	}

	// 获取列信息
	columns, err := d.getColumns(tableName)
	if err != nil {
		return nil, err
	}

	// 获取主键信息
	primaryKeys, err := d.getPrimaryKeys(tableName)
	if err != nil {
		return nil, err
	}

	// 标记主键列
	for i := range columns {
		for _, pk := range primaryKeys {
			if columns[i].Name == pk {
				columns[i].IsPrimaryKey = true
				break
			}
		}
	}

	return &TableInfo{
		Name:        tableName,
		Comment:     tableComment,
		Columns:     columns,
		PrimaryKeys: primaryKeys,
	}, nil
}

// getTableComment 获取表注释
func (d *Database) getTableComment(tableName string) (string, error) {
	query := `
		SELECT table_comment 
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name = ?
	`
	var comment string
	err := d.db.QueryRow(query, tableName).Scan(&comment)
	if err != nil {
		return "", fmt.Errorf("获取表注释失败: %v", err)
	}
	return comment, nil
}

// getColumns 获取列信息
func (d *Database) getColumns(tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			column_name,
			column_type,
			is_nullable,
			column_default,
			column_comment,
			extra
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY ordinal_position
	`

	rows, err := d.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("查询列信息失败: %v", err)
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var extra string
		var isNullableStr string
		var defaultVal sql.NullString
		var comment sql.NullString
		err := rows.Scan(
			&col.Name,
			&col.Type,
			&isNullableStr,
			&defaultVal,
			&comment,
			&extra,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描列信息失败: %v", err)
		}

		// 转换字符串为布尔值
		col.IsNullable = strings.ToUpper(isNullableStr) == "YES"

		// 处理默认值
		if defaultVal.Valid {
			col.DefaultValue = defaultVal.String
		} else {
			col.DefaultValue = ""
		}

		// 处理注释
		if comment.Valid {
			col.Comment = comment.String
		} else {
			col.Comment = ""
		}

		// 检查是否为自增列
		col.IsAutoIncr = strings.Contains(strings.ToLower(extra), "auto_increment")

		columns = append(columns, col)
	}

	return columns, nil
}

// getPrimaryKeys 获取主键信息
func (d *Database) getPrimaryKeys(tableName string) ([]string, error) {
	query := `
		SELECT column_name
		FROM information_schema.key_column_usage
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND constraint_name = 'PRIMARY'
		ORDER BY ordinal_position
	`

	rows, err := d.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("查询主键信息失败: %v", err)
	}
	defer rows.Close()

	var primaryKeys []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, fmt.Errorf("扫描主键信息失败: %v", err)
		}
		primaryKeys = append(primaryKeys, columnName)
	}

	return primaryKeys, nil
}

// ListTables 列出所有表
func (d *Database) ListTables() ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("扫描表名失败: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
