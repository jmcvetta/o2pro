// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"strings"
)

type AuthTemplate struct {
	Username string
	Scopes   []string
	Note     string
}

// ScopesMap returns a map of the scopes in this authorization, for easy look
// up.  Bool is always true.
func (a *Authz) ScopesMap() map[string]bool {
	sm := map[string]bool{}
	for _, s := range a.Scopes {
		sm[s] = true
	}
	return sm
}

func (a *Authz) ScopeString() string {
	return strings.Join(a.Scopes, " ")
}
