// Copyright (C) 2013 Jason McVetta, all rights reserved.

package o2pro

import "github.com/coocood/qbs"

type qbsStorage struct {
	q      *qbs.Qbs
	dbName string // Consumbed by qbs migrations
}

// MigrateTables will attempt to migrate the database to the current schema,
// creating tables that do not exist, and adding columns to those that do.
// Only additive operations are supported - it will not alter or delete columns
// - so it should be safe for production. Will panic if it can't migrate a
// table, or return an error if it cannot create a necessary index.
func (q *qbsStorage) Migrate() error {
	m := qbs.NewMigration(q.q.Db, q.dbName, q.q.Dialect)
	for _, t := range tables {
		err := m.CreateTableIfNotExists(t)
		if err != nil {
			return err
		}

	}
	return nil
}

// dropTables drops all tables.  It can ONLY be used on databases whose
// names end in "_test".
func (q *qbsStorage) dropTables() error {
	m := qbs.NewMigration(q.q.Db, q.dbName, q.q.Dialect)
	for _, t := range tables {
		m.DropTable(t)
	}
	return nil
}
