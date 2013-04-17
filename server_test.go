// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"github.com/bmizerany/assert"
	"labix.org/v2/mgo"
	"log"
	"net/url"
	"testing"
	"time"
)

var (
	testScopesAll     = []string{"enterprise", "shuttlecraft", "intrepid"}
	testScopesDefault = []string{"shuttlecraft"}
)

// fakeAuth authorizes everyone for everything.
func fakeAuth(*url.Userinfo, AuthRequest) (bool, error) {
	return true, nil
}

func col(db *mgo.Database) *mgo.Collection {
	return db.C("authorizations")
}

func setup(t *testing.T) (*Server, *mgo.Database) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	db := session.DB("test_btoken")
	s, err := NewMongoServer(db, DefaultExpireAfter, fakeAuth)
	if err != nil {
		t.Fatal(err)
	}
	s.Scopes = testScopesAll
	s.DefaultScopes = testScopesDefault
	return s, db
}

func TestNewAuth(t *testing.T) {
	s, db := setup(t)
	owner := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		Scopes: scopes,
	}
	auth, err := s.NewAuth(owner, req)
	if err != nil {
		t.Error(err)
	}
	c := col(db)
	query := struct {
		Token string
	}{
		Token: auth.Token,
	}
	q := c.Find(&query)
	cnt, err := q.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
	a := Authorization{}
	err = q.One(&a)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, owner, a.Owner)
	for _, scope := range scopes {
		_, ok := a.ScopesMap[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func TestGetAuthorization(t *testing.T) {
	s, _ := setup(t)
	owner := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		Scopes: scopes,
	}
	auth, err := s.NewAuth(owner, req)
	if err != nil {
		t.Error(err)
	}
	a, err := s.GetAuth(auth.Token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, owner, a.Owner)
	for _, scope := range scopes {
		_, ok := a.ScopesMap[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func TestExpiration(t *testing.T) {
	five, _ := time.ParseDuration("5ms")
	seven, _ := time.ParseDuration("7ms")
	s, _ := setup(t)
	owner := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		Scopes:      scopes,
		ExpireAfter: five,
	}
	auth, err := s.NewAuth(owner, req)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(seven) // Authorization should be expired
	_, err = s.GetAuth(auth.Token)
	assert.Equal(t, ErrInvalidToken, err)
}
