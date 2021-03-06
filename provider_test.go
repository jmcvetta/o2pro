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

func doTestAuthz(p *Provider, t *testing.T) {
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	auth, err := p.NewAuthz(username, "", scopes)
	if err != nil {
		t.Error(err)
	}
	a, err := p.Authz(auth.Token)
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

func doTestExpiration(p *Provider, t *testing.T) {
	five, _ := time.ParseDuration("5ms")
	seven, _ := time.ParseDuration("7ms")
	p.Duration = five
	username := "jtkirk"
	scopes := []string{"enterprise", "shuttlecraft"}
	auth, err := p.NewAuthz(username, "", scopes)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(seven) // Authz should be expired
	_, err = p.Authz(auth.Token)
	if err != ErrInvalidToken {
		t.Fatal(err)
	}
}

// for testing things that do not depend on storage
type nullStorage struct {
}

func (n *nullStorage) saveAuthz(a *Authz) error {
	return ErrNotImplemented
}

func (n *nullStorage) authz(token string) (*Authz, error) {
	return nil, ErrNotImplemented
}

func (n *nullStorage) initialize() error {
	return ErrNotImplemented
}

func (n *nullStorage) migrate() error {
	return ErrNotImplemented
}

// for testing things that do not depend on storage
func testNull(t *testing.T) *Provider {
	p := NewProvider(&nullStorage{}, kirkAuthenticator, GrantAll)
	p.Scopes = testScopesAll
	p.DefaultScopes = testScopesDefault
	return p
}

func TestPrettyPrint(t *testing.T) {
	prettyPrint(testScopesAll)
}
