// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import "net/http"

/*
ACCESSING PROTECTED RESOURCES
http://tools.ietf.org/html/rfc6749#section-7
*/

import ()

type AccessController struct {
}

func (c *AccessController) Authz(token string) (*Authz, error) {
	return nil, ErrNotImplemented
}

func (c *AccessController) ReqAuthz(r *http.Request) (*Authz, error) {
	return nil, ErrNotImplemented
}

// ProtectScope wraps a HandlerFunc, restricting access to authenticated users
// with the specified scope.
func (c *AccessController) ProtectScope(fn http.HandlerFunc, scope string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
