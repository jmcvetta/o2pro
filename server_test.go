// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

var (
	testScopesAll     = []string{"enterprise", "shuttlecraft", "intrepid"}
	testScopesDefault = []string{"shuttlecraft"}
)

// An Authenticator implementation that authenticates user "jtkirk" with
// password "Beam me up, Scotty!".
func kirkAuthenticator(username, password string) (bool, error) {
	if username == "jtkirk" && password == "Beam me up, Scotty!" {
		return true, nil
	}
	return false, nil
}

func doTestAuthz(s *Server, t *testing.T) {
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	tmpl := AuthzTemplate{
		User:   username,
		Scopes: scopes,
	}
	auth, err := s.NewAuthz(tmpl)
	if err != nil {
		t.Error(err)
	}
	a, err := s.Authz(auth.Token)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, username, a.User)
	sm := a.ScopesMap()
	for _, scope := range scopes {
		_, ok := sm[scope]
		assert.T(t, ok, "Expected scope: ", scope)
	}
}

func doTestExpiration(s *Server, t *testing.T) {
	five, _ := time.ParseDuration("5ms")
	seven, _ := time.ParseDuration("7ms")
	s.Duration = five
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	tmpl := AuthzTemplate{
		User:   username,
		Scopes: scopes,
	}
	auth, err := s.NewAuthz(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(seven) // Authz should be expired
	_, err = s.Authz(auth.Token)
	if err != ErrInvalidToken {
		t.Fatal(err)
	}
}

// for testing things that do not depend on storage
type nullStorage struct {
}

func (n *nullStorage) SaveAuthz(a *Authz) error {
	return nil
}

func (n *nullStorage) Authz(token string) (*Authz, error) {
	return nil, nil
}

func (n *nullStorage) Initialize() error {
	return nil
}

func (n *nullStorage) Migrate() error {
	return nil
}

// for testing things that do not depend on storage
func testNull(t *testing.T) *Server {
	s := NewServer(&nullStorage{}, kirkAuthenticator, GrantAll)
	s.Scopes = testScopesAll
	s.DefaultScopes = testScopesDefault
	return s
}

func TestPrettyPrint(t *testing.T) {
	prettyPrint(testScopesAll)
}
