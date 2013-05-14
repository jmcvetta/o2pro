// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/restclient"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func doTestPasswordRequest(s *Server, t *testing.T) {
	//
	// Prepare handler
	//
	h := s.HandlerFunc(PasswordGrant)
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	//
	// REST request
	//
	scopes := []string{"enterprise", "intrepid"}
	scopeStr := strings.Join(scopes, " ")
	username := "jtkirk"
	password := "Beam me up, Scotty!"
	u := url.UserPassword(username, password)
	preq := PasswordRequest{
		GrantType: "password",
		Username:  "jtkirk",
		Password:  password,
		Scope:     scopeStr,
		Note:      "foo bar baz",
	}
	var res TokenResponse
	var e interface{}
	rr := restclient.RequestResponse{
		Url:      hserv.URL,
		Method:   "POST",
		Userinfo: u,
		Data:     &preq,
		Result:   &res,
		Error:    &e,
	}
	c := restclient.New()
	c.UnsafeBasicAuth = true
	status, err := c.Do(&rr)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, status)
	assert.NotEqual(t, nil, uuid.Parse(res.AccessToken))
	assert.Equal(t, scopeStr, res.Scope)
	assert.Equal(t, "bearer", res.TokenType)
}

func TestPasswordBadCreds(t *testing.T) {
	s := testNull(t)
	//
	// Prepare handler
	//
	h := s.HandlerFunc(PasswordGrant)
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	//
	// REST request
	//
	username := "jtkirk"
	password := "Go Klingons!"
	u := url.UserPassword(username, password)
	preq := PasswordRequest{
		GrantType: "password",
		Username:  "jtkirk",
		Password:  password,
	}
	var res string
	var e interface{}
	rr := restclient.RequestResponse{
		Url:      hserv.URL,
		Method:   "POST",
		Userinfo: u,
		Data:     &preq,
		Result:   &res,
		Error:    &e,
	}
	c := restclient.New()
	c.UnsafeBasicAuth = true
	status, err := c.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
}

func TestPasswordBadAuthHeader(t *testing.T) {
	s := testNull(t)
	//
	// Prepare handler
	//
	h := s.HandlerFunc(PasswordGrant)
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	//
	// Regex doesn't match
	//
	req, err := http.NewRequest("POST", hserv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "foobar")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, resp.StatusCode)
	//
	// Base64 decode failed
	//
	req, err = http.NewRequest("POST", hserv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic foobar")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, resp.StatusCode)
	//
	// String split failed
	//
	req, err = http.NewRequest("POST", hserv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	str := base64.URLEncoding.EncodeToString([]byte("foobar"))
	req.Header.Add("Authorization", "Basic "+str)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 400, resp.StatusCode)
}
