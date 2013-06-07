// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

/*
Implementation of RESOURCE OWNER PASSWORD CREDENTIALS GRANT workflow.
http://tools.ietf.org/html/rfc6749#section-4.3
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// A PasswordRequest is submitted by a client requesting authorization using
// the Resource Owner Password Credentials Grant flow.
type PasswordRequest struct {
	GrantType string `json:"grant_type"` // REQUIRED.  Value MUST be set to "password".
	Username  string `json:"username"`   // REQUIRED.  The resource owner username.
	Password  string `json:"password"`   // REQUIRED.  The resource owner password.
	Scope     string `json:"scope"`      // OPTIONAL.  The scope of the access request as described by http://tools.ietf.org/html/rfc6749#section-3.3
	Note      string `json:"note"`       // OPTIONAL.  Not part of RFC spec - inspired by Github.
}

// PasswordGrant supports authorization via the  Resource Owner Password
// Credentials Grant workflow.
func passwordGrant(p *Provider, w http.ResponseWriter, r *http.Request) {
	//
	// Authenticate
	//
	malformed := "Malformed Authorization header"
	username, password, err := BasicAuth(r)
	if err != nil {
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	ok, err := p.Authenticate(username, password)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Invalid username/password", http.StatusUnauthorized)
		return
	}
	//
	// Parse authorization request
	//
	dec := json.NewDecoder(r.Body)
	var preq PasswordRequest
	err = dec.Decode(&preq)
	if err != nil && err.Error() != "EOF" {
		log.Println(err)
		msg := "Missing or bad request body"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if username != preq.Username || preq.GrantType != "password" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	//
	// Validate scope
	//
	scopes := strings.Split(preq.Scope, " ")
	valid := sliceMap(p.Scopes)
	for _, scope := range scopes {
		_, ok = valid[scope]
		if !ok {
			http.Error(w, "Invalid scope: "+scope, http.StatusBadRequest)
			return
		}
		ok, err = p.Grant(username, scope, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Not authorized for scope: "+scope, http.StatusUnauthorized)
			return
		}
	}
	//
	// Create new authorization
	//
	a, err := p.NewAuthz(preq.Username, preq.Note, scopes)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	//
	// Compose response
	//
	resp := TokenResponse{
		AccessToken: a.Token,
		TokenType:   "bearer",
		ExpiresIn:   int(a.Expiration.Sub(time.Now()).Seconds()),
		Scope:       a.ScopeString(),
	}
	enc := json.NewEncoder(w)
	enc.Encode(&resp)
	return
}
