// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/restclient"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPasswordRequest(t *testing.T) {
	//
	// Prepare handler
	//
	s, _ := setup(t)
	h := s.PasswordHandler()
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
