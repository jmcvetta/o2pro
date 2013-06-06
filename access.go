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

// RequireScope wraps a HandlerFunc, restricting access to authenticated users
// with the specified scope.
func (p *Provider) RequireScope(fn http.HandlerFunc, scope string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := BearerToken(r)
		if err != nil { // No token found
			log.Println(err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		a, err := p.authz(token)
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

// RequireAuthc wraps a HandlerFunc, restricting access to authenticated users.
func (p *Provider) RequireAuthc(fn http.HandlerFunc, scope string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := BearerToken(r)
		if err != nil { // No token found
			log.Println(err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		_, err = p.authz(token)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		fn(w, r)
		return
	}
}
