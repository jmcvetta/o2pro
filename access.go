// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import "log"

import "net/http"

/*
ACCESSING PROTECTED RESOURCES
http://tools.ietf.org/html/rfc6749#section-7
*/

import ()

// An AccessController restricts access to resources using OAuth tokens.
type AccessController struct {
	Storage
}

// ProtectScope wraps a HandlerFunc, restricting access to authenticated users
// with the specified scope.
func (c *AccessController) ProtectScope(fn http.HandlerFunc, scope string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := BearerToken(r)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		a, err := c.authz(token)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		_, ok := a.ScopesMap()[scope]
		if !ok {
			log.Printf("Need scope '%v' but only authorized for '%v'", scope, a.ScopeString())
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		fn(w, r) // Call the wrapped function
		return
	}
}

// running inside the same program.
func NewAccessController(s Storage) *AccessController {
	return &AccessController{s}
}

func (c *AccessController) ReqAuthz(r *http.Request) (*Authz, error) {
	return nil, ErrNotImplemented
}
