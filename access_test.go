// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"github.com/bmizerany/assert"
	"github.com/jmcvetta/restclient"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func doTestAccessController(p *Provider, t *testing.T) {
	c := NewAccessController(p)
	h := c.ProtectScope(fooHandler, "enterprise")
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	note := "foo bar baz"
	auth, _ := p.NewAuthz(username, note, scopes)
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+auth.Token)
	rr := restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, status)
}

func doTestAccessControllerBadScope(p *Provider, t *testing.T) {
	c := NewAccessController(p)
	h := c.ProtectScope(fooHandler, "foobar") // Not among the authorized scopes
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	note := "foo bar baz"
	auth, _ := p.NewAuthz(username, note, scopes)
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+auth.Token)
	rr := restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
		Header: &header,
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
}

func doTestAccessControllerNoToken(p *Provider, t *testing.T) {
	c := NewAccessController(p)
	h := c.ProtectScope(fooHandler, "enterprise")
	hserv := httptest.NewServer(h)
	defer hserv.Close()
	rr := restclient.RequestResponse{
		Url:    hserv.URL,
		Method: "GET",
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 401, status)
}
