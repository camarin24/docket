package docket

import (
	"github.com/camarin24/docket/internal/db"
	"github.com/camarin24/docket/pkg/types"
	"gorm.io/gorm"
)

type DbAdapter interface {
	NewDb() *Db
}

type SqliteDb struct {
	DbName string
}

type PostgresDb struct {
	DbHost     string
	DbName     string
	DbUser     string
	DbPassword string
	SSLMode    string
	DbPort     string
}

type Db struct {
	*gorm.DB
}

const (
	DefaultSqliteDbName    = "docket.sqlite3"
	DefaultPostgresSSLMode = "disabled"
	DefaultPostgresPort    = "5432"
)

func (sq SqliteDb) NewDb() *Db {
	if sq.DbName == "" {
		sq.DbName = DefaultSqliteDbName
	}

	return &Db{
		DB: db.NewSqlite(db.DbConfig{
			DbName: sq.DbName,
		}),
	}
}

func (pg PostgresDb) NewDb() *Db {
	if pg.DbPort == "" {
		pg.DbPort = DefaultPostgresPort
	}

	if pg.SSLMode == "" {
		pg.SSLMode = DefaultPostgresSSLMode
	}

	return &Db{db.NewPostgres(db.DbConfig{
		DbHost: pg.DbHost, DbName: pg.DbName, DbUser: pg.DbUser, DbPassword: pg.DbPassword, DbPort: pg.DbPort, SSLMode: pg.SSLMode,
	})}
}

func (db *Db) GetDocumentsNameByStorageKey(storageKey string) []types.Document {
	var docs []types.Document
	db.Where(&types.Document{StorageKey: storageKey}).Select("name").Find(&docs)
	return docs
}

func (db *Db) CreateDocuments(docs ...types.Document) {
	db.Create(&docs)
}

func (db *Db) Migrate() {
	err := db.AutoMigrate(types.Document{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(types.Metadata{})
	if err != nil {
		panic(err)
	}
}