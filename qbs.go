// Copyright (C) 2013 Jason McVetta, all rights reserved.

package o2pro

import (
	"database/sql"
	"github.com/coocood/qbs"
	"log"
)

// tables is the list of structs defining tables for qbs to migrate.
var tables = []interface{}{
	&Client{},
	&Authz{},
	&AuthzScope{},
	&Code{},
}

type qbsStorage struct {
	driver  string
	dsn     string
	dialect qbs.Dialect
}

func NewQbsStorage(driverName, dataSourceName string) Storage {
	q := &qbsStorage{
		driver: driverName,
		dsn:    dataSourceName,
	}
	switch driverName {
	default:
		log.Panic("Invalid driverName")
	case "postgres":
		q.dialect = qbs.NewPostgres()
	case "mysql":
		q.dialect = qbs.NewMysql()
	}
	return q
}

func (s *qbsStorage) SaveAuthz(a *Authz) error {
	q := s.qbs()
	defer q.Close()
	_, err := q.Save(a)
	return err
}

func (s *qbsStorage) Authz(token string) (*Authz, error) {
	q := s.qbs()
	defer q.Close()
	a := &Authz{
		Token: token,
	}
	err := q.Find(a)
	if err == sql.ErrNoRows {
		return a, ErrInvalidToken
	}
	return a, err
}

func (s *qbsStorage) Initialize() error {
	return s.Migrate()
}

// MigrateTables will attempt to migrate the database to the current schema,
// creating tables that do not exist, and adding columns to those that do.
// Only additive operations are supported - it will not alter or delete columns
// - so it should be safe for production. Will panic if it can't migrate a
// table, or return an error if it cannot create a necessary index.
func (q *qbsStorage) Migrate() error {
	m := qbs.NewMigration(q.qbs().Db, "", q.qbs().Dialect)
	for _, t := range tables {
		err := m.CreateTableIfNotExists(t)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *qbsStorage) qbs() *qbs.Qbs {
	var err error
	db := qbs.GetFreeDB()
	if db == nil {
		db, err = sql.Open(s.driver, s.dsn)
		if err != nil {
			log.Panic(err)
		}
	}
	return qbs.New(db, s.dialect)
}
