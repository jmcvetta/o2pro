// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package btoken

import (
	"github.com/bmizerany/assert"
	"labix.org/v2/mgo"
	"testing"
)

var (
	testAS    AuthServer
	testMongo *mgo.Database
)

func col() *mgo.Collection {
	return testMongo.C("authorizations")
}

func setup(t *testing.T) {
	/*
		if testAS != nil {
			t.Log("Using existing testAuthServer\n")
			return
		}
		t.Log("Initializing testAuthServer\n")
	*/
	session, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	testMongo = session.DB("test_btoken")
	err = testMongo.DropDatabase()
	if err != nil {
		t.Fatal(err)
	}
	testAS, err = NewMongoAuthServer(testMongo)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestIssueToken(t *testing.T) {
	setup(t)
	user := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	req := AuthRequest{
		User:   user,
		Scopes: scopes,
	}
	token, err := testAS.IssueToken(req)
	if err != nil {
		t.Error(err)
	}
	c := col()
	cnt, err := c.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
	query := struct {
		Token string
	}{
		Token: token,
	}
	q := c.Find(query)
	cnt, err = q.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
	a := Authorization{}
	err = q.One(&a)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, user, a.User)
	for _, scope := range scopes {
		_, ok := a.Scopes[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func TestFoobar(t *testing.T) {
	setup(t)

}
