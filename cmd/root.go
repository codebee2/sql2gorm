package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"sql2gorm/config"
	"sql2gorm/database"
	"sql2gorm/generator"
	"sql2gorm/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	sqlFile     string
	outputFile  string
	tableName   string
	packageName string
	dsn         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sql2gorm",
	Short: "SQL表结构转GORM模型工具",
	Long: `SQL2GORM是一个将MySQL表结构SQL语句转换为GORM模型代码的工具。

支持通过配置文件自定义生成规则，自动生成符合GORM规范的Go结构体代码。

主要特性:
- 自动解析CREATE TABLE SQL语句
- 智能生成GORM模型结构体
- 自动添加JSON和GORM标签
- 支持自定义类型映射规则
- 批量生成到指定文件`,
	RunE: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: config.yaml)")
	rootCmd.PersistentFlags().StringVar(&sqlFile, "sql", "", "SQL文件路径")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "输出Go文件路径")
	rootCmd.PersistentFlags().StringVar(&tableName, "table", "", "表名")
	rootCmd.PersistentFlags().StringVar(&packageName, "package", "models", "Go包名")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "数据库连接字符串 (DSN)")

	// 绑定viper
	viper.BindPFlag("sql_file", rootCmd.PersistentFlags().Lookup("sql"))
	viper.BindPFlag("output_file", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("table_name", rootCmd.PersistentFlags().Lookup("table"))
	viper.BindPFlag("package_name", rootCmd.PersistentFlags().Lookup("package"))
	viper.BindPFlag("dsn", rootCmd.PersistentFlags().Lookup("dsn"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 搜索配置文件
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 如果找到配置文件，读取它
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) error {
	// 加载配置
	cfg := &config.Config{
		SQLFile:     viper.GetString("sql_file"),
		OutputFile:  viper.GetString("output_file"),
		TableName:   viper.GetString("table_name"),
		PackageName: viper.GetString("package_name"),
		DSN:         viper.GetString("dsn"),
	}

	// 验证配置
	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	var tableInfo *parser.TableInfo
	var err error

	// 优先使用DSN方式
	if cfg.DSN != "" {
		// 从数据库读取表结构
		tableInfo, err = getTableInfoFromDB(cfg.DSN, cfg.TableName)
		if err != nil {
			return fmt.Errorf("从数据库读取表结构失败: %v", err)
		}
		fmt.Printf("✅ 成功从数据库读取表结构: %s\n", cfg.TableName)
	} else if cfg.SQLFile != "" {
		// 从SQL文件读取
		tableInfo, err = getTableInfoFromSQL(cfg.SQLFile)
		if err != nil {
			return fmt.Errorf("从SQL文件读取表结构失败: %v", err)
		}
		fmt.Printf("✅ 成功从SQL文件读取表结构: %s\n", cfg.SQLFile)
	} else {
		return fmt.Errorf("必须提供DSN或SQL文件路径")
	}

	// 生成GORM模型
	generator := generator.NewGORMGenerator(cfg)
	modelCode, err := generator.Generate(tableInfo)
	if err != nil {
		return fmt.Errorf("生成GORM模型失败: %v", err)
	}

	// 写入输出文件
	if err := os.WriteFile(cfg.OutputFile, []byte(modelCode), 0644); err != nil {
		return fmt.Errorf("写入输出文件失败: %v", err)
	}

	fmt.Printf("✅ 成功生成GORM模型到文件: %s\n", cfg.OutputFile)
	return nil
}

// getTableInfoFromDB 从数据库读取表结构
func getTableInfoFromDB(dsn, tableName string) (*parser.TableInfo, error) {
	db, err := database.NewDatabase(dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tableInfo, err := db.GetTableInfo(tableName)
	if err != nil {
		return nil, err
	}

	return tableInfo.ToParserTableInfo(), nil
}

// getTableInfoFromSQL 从SQL文件读取表结构
func getTableInfoFromSQL(sqlFile string) (*parser.TableInfo, error) {
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return nil, fmt.Errorf("读取SQL文件失败: %v", err)
	}

	parser := parser.NewSQLParser()
	tableInfo, err := parser.Parse(string(sqlContent))
	if err != nil {
		return nil, fmt.Errorf("解析SQL失败: %v", err)
	}

	return tableInfo, nil
}

func validateConfig(cfg *config.Config) error {
	if cfg.OutputFile == "" {
		return fmt.Errorf("输出文件路径不能为空")
	}
	if cfg.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}

	// 检查是否提供了DSN或SQL文件
	if cfg.DSN == "" && cfg.SQLFile == "" {
		return fmt.Errorf("必须提供DSN或SQL文件路径")
	}

	// 如果使用SQL文件方式，检查文件是否存在
	if cfg.SQLFile != "" {
		if _, err := os.Stat(cfg.SQLFile); os.IsNotExist(err) {
			return fmt.Errorf("SQL文件不存在: %s", cfg.SQLFile)
		}
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(cfg.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	return nil
}
