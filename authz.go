// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

// An Authz is an authorization.
type Authz struct {
	Id         int64 `bson:",omitempty`
	Uuid       string
	Token      string
	User       string
	ClientId   int64
	Client     *Client
	Issued     time.Time
	Expiration time.Time
	Note       string
	Scopes     []string
}

// NewAuth issues a new authorization.
func (p *Provider) NewAuthz(user, note string, scopes []string) (*Authz, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := uuid.New() + string(b)
	token = base64.StdEncoding.EncodeToString([]byte(token))
	a := Authz{
		Token:      token,
		Uuid:       uuid.New(),
		User:       user,
		Scopes:     scopes,
		Expiration: time.Now().Add(p.Duration),
		Note:       note,
	}
	err = p.SaveAuthz(&a)
	return &a, err
}

// SaveAuthz saves an authorization to storage.
func (p *Provider) SaveAuthz(a *Authz) error {
	return p.saveAuthz(a)
}

// Authz looks up an authorization based on its token.
func (p *Provider) Authz(token string) (*Authz, error) {
	return p.authz(token)
}

// ScopesMap returns a map of the scopes in this authorization, for easy look
// up.  Bool is always true.
func (a *Authz) ScopesMap() map[string]bool {
	return sliceMap(a.Scopes)
}

func (a *Authz) ScopeString() string {
	return strings.Join(a.Scopes, " ")
}
