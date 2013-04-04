// Copyright (C) 2013 Jason McVetta, all rights reserved.

package btoken

import (
	"labix.org/v2/mgo"
	"github.com/bmizerany/assert"
	"testing"
)

var (
	testAS AuthServer
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
	req := AuthRequest{
		User:   "jtkirk",
		Scopes: []string{"enterprise", "shuttlecraft"},
	}
	token, err := testAS.IssueToken(req)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Token: %v\n", token)
	c := col()
	cnt, err := c.Count()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, cnt)
}

func TestFoobar(t *testing.T) {
	setup(t)

}
