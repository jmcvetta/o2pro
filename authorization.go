// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

type AuthTemplate struct {
	Username string
	Scopes   []string
	Note     string
}

type Auth struct {
	Id         bson.ObjectId `bson:"_id",json:"id"` // Unique storage-dependent ID for this Authorization
	Token      string        `json:"token"`
	Username   string        `json:"username"`
	Scopes     []string      `json:"scopes"`
	Expiration time.Time     `json:"expiration"`
	Note       string        `json:"note"`
}

// ScopesMap returns a map of the scopes in this authorization, for easy look
// up.  Bool is always true.
func (a *Auth) ScopesMap() map[string]bool {
	sm := map[string]bool{}
	for _, s := range a.Scopes {
		sm[s] = true
	}
	return sm
}

func (a *Auth) ScopeString() string {
	return strings.Join(a.Scopes, " ")
}
