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

// NewAuth issues a new authorization.
func (s *Server) NewAuthz(user, note string, scopes []string) (*Authz, error) {
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
		Expiration: time.Now().Add(s.Duration),
		Note:       note,
	}
	err = s.SaveAuthz(&a)
	return &a, err
}

// SaveAuthz saves an authorization to storage.
func (s *Server) SaveAuthz(a *Authz) error {
	return s.saveAuthz(a)
}

// Authz looks up an authorization based on its token.
func (s *Server) Authz(token string) (*Authz, error) {
	return s.authz(token)
}

// ScopesMap returns a map of the scopes in this authorization, for easy look
// up.  Bool is always true.
func (a *Authz) ScopesMap() map[string]bool {
	return sliceMap(a.Scopes)
}

func (a *Authz) ScopeString() string {
	return strings.Join(a.Scopes, " ")
}
