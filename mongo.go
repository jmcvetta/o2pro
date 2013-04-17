// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// NewMongoAuthServer configures a MongoDB-based AuthServer.  If expireAfter is
// not nil, authorizations will be automatically expired.
func NewMongoServer(db *mgo.Database, duration string, a Authorizer) (*Server, error) {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	stor := mongoStorage{
		db:          db,
		name:        "authorizations",
		expireAfter: dur,
	}
	serv := Server{
		Storage:       &stor,
		MaxDuration:   dur,
		Authorizer:    a,
		Logger:        DefaultLogger,
		Scopes:        DefaultScopes,
		DefaultScopes: DefaultScopes,
	}
	err = serv.Activate()
	return &serv, err
}

type mongoStorage struct {
	db          *mgo.Database
	name        string // Collection name
	expireAfter time.Duration
}

func (m *mongoStorage) Activate() error {
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

func (s *mongoStorage) GetAuth(token string) (Authorization, error) {
	a := Authorization{}
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

func (s *mongoStorage) SaveAuth(auth *Authorization) error {
	oid := bson.NewObjectId()
	auth.AuthId = oid
	return s.col().Insert(auth)
}

// col returns a Collection object in a new mgo session
func (s *mongoStorage) col() *mgo.Collection {
	session := s.db.Session.Copy()
	d := session.DB(s.db.Name)
	return d.C(s.name)
}
