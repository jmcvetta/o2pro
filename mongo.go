// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"labix.org/v2/mgo"
	"time"
)

// func NewMongoStorage(db *mgo.Database, dur time.Duration, a Authorizer) (*Storage, error) {
func NewMongoStorage(db *mgo.Database, dur time.Duration) Storage {
	return &mongoStorage{
		db:          db,
		name:        "authorizations",
		expireAfter: dur,
	}
}

type mongoStorage struct {
	db          *mgo.Database
	name        string // Collection name
	expireAfter time.Duration
}

func (m *mongoStorage) Initialize() error {
	return m.Migrate()
}

func (m *mongoStorage) Migrate() error {
	//
	// Declare Indexes
	//
	idxs := []mgo.Index{
		mgo.Index{
			Key:      []string{"token"},
			Unique:   true,
			DropDups: false,
		},
		mgo.Index{
			Key:      []string{"expiration"},
			Unique:   true,
			DropDups: false,
		},
	}
	c := m.col()
	for _, i := range idxs {
		err := c.EnsureIndex(i)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *mongoStorage) Authz(token string) (*Authz, error) {
	a := new(Authz)
	c := s.col()
	query := struct {
		Token string
	}{
		Token: token,
	}
	q := c.Find(query)
	cnt, err := q.Count()
	if err != nil {
		return a, err
	}
	if cnt < 1 {
		return a, ErrInvalidToken
	}
	err = q.One(&a)
	if err != nil {
		return a, err
	}
	if time.Now().After(a.Expiration) {
		c.Remove(query)
		return a, ErrInvalidToken
	}
	return a, nil
}

func (s *mongoStorage) SaveAuthz(a *Authz) error {
	return s.col().Insert(a)
}

// col returns a Collection object in a new mgo session
func (s *mongoStorage) col() *mgo.Collection {
	session := s.db.Session.Copy()
	d := session.DB(s.db.Name)
	return d.C(s.name)
}
