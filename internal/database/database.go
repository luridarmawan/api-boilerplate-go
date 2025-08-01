package database

import (
	"fmt"
	"log"

	"apiserver/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var tablePrefix string

// CustomNamingStrategy implements schema.Namer interface
type CustomNamingStrategy struct {
	schema.NamingStrategy
	TablePrefix string
}

// JoinTableName returns the join table name with prefix
func (ns CustomNamingStrategy) JoinTableName(joinTable string) string {
	if ns.TablePrefix == "" {
		return joinTable
	}
	return ns.TablePrefix + joinTable
}

// TableName returns the table name with prefix
func (ns CustomNamingStrategy) TableName(table string) string {
	if ns.TablePrefix == "" {
		return ns.NamingStrategy.TableName(table)
	}
	return ns.TablePrefix + ns.NamingStrategy.TableName(table)
}

func InitDatabase(config *configs.Config) {
	// Set the global table prefix
	tablePrefix = config.DBTablePrefix
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		config.DBSSLMode,
	)

	// Create custom naming strategy with prefix
	namingStrategy := CustomNamingStrategy{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		TablePrefix: config.DBTablePrefix,
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}

// GetTableName returns the table name with prefix
func GetTableName(tableName string) string {
	if tablePrefix == "" {
		return tableName
	}
	return tablePrefix + tableName
}

// GetTablePrefix returns the current table prefix
func GetTablePrefix() string {
	return tablePrefix
}
