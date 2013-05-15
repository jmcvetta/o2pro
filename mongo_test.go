// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"github.com/bmizerany/assert"
	"labix.org/v2/mgo"
	"log"
	"testing"
	"time"
)

func col(db *mgo.Database) *mgo.Collection {
	return db.C("authorizations")
}

func testMongo(t *testing.T) (*Provider, *mgo.Database) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	db := session.DB("test_o2pro")
	dur, err := time.ParseDuration(DefaultExpireAfter)
	if err != nil {
		t.Fatal(err)
	}
	stor := NewMongoStorage(db, dur)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProvider(stor, kirkAuthenticator, GrantAll)
	p.Scopes = testScopesAll
	p.DefaultScopes = testScopesDefault
	err = p.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	err = p.Migrate()
	if err != nil {
		t.Fatal(err)
	}
	return p, db
}

func TestMgoNewAuth(t *testing.T) {
	s, db := testMongo(t)
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	note := "foo bar baz"
	auth, err := s.NewAuthz(username, note, scopes)
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
	a := Authz{}
	err = q.One(&a)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, username, a.User)
	sm := a.ScopesMap()
	for _, scope := range scopes {
		_, ok := sm[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func TestMgoAuthz(t *testing.T) {
	s, _ := testMongo(t)
	doTestAuthz(s, t)
}

func TestMgoExpiration(t *testing.T) {
	s, _ := testMongo(t)
	doTestExpiration(s, t)
}

func TestMgoPasswordRequest(t *testing.T) {
	s, _ := testMongo(t)
	doTestPasswordRequest(s, t)
}
