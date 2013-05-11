// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"github.com/coocood/qbs"
	"github.com/darkhelmet/env"
	"github.com/lib/pq"
	"log"
	"testing"
)

func testQbs(t *testing.T) (*Server, *qbs.Qbs) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// PostgreSQL & Qbs
	//
	dbUrl := env.StringDefault("DATABASE_URL", "postgres://")
	dsn, err := pq.ParseURL(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	stor := NewQbsStorage("postgres", dsn)
	err = stor.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	qs := stor.(*qbsStorage)
	q := qs.qbs()
	// s, err := NewMongoServer(db, DefaultExpireAfter, kirkAuthorizer)
	s := NewServer(stor, kirkAuthorizer)
	if err != nil {
		t.Fatal(err)
	}
	s.Scopes = testScopesAll
	s.DefaultScopes = testScopesDefault
	return s, q
}

func TestQbsNewAuth(t *testing.T) {
	s, q := testQbs(t)
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	tmpl := AuthTemplate{
		User:   username,
		Scopes: scopes,
		Note:   "foo bar baz",
	}
	a0, err := s.NewAuthz(tmpl)
	if err != nil {
		t.Error(err)
	}
	prettyPrint(a0)
	a1 := new(Authz)
	a1.Id = a0.Id
	// aSlice := []*Authz{&a1,}
	// prettyPrint(aSlice)
	// err = q.FindAll(&aSlice)
	err = q.Find(a1)
	q.Query()
	if err != nil {
		t.Error(err)
	}
	prettyPrint(a1)
	/*
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
	*/
}

func TestQbsAuthz(t *testing.T) {
	s, _ := testQbs(t)
	doTestAuthz(s, t)
}

func TestQbsExpiration(t *testing.T) {
	s, _ := testQbs(t)
	doTestExpiration(s, t)
}

func TestQbsPasswordRequest(t *testing.T) {
	s, _ := testQbs(t)
	doTestPasswordRequest(s, t)
}
