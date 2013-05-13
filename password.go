// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

/*
	Implementation of the the Resource Owner Password Credentials Grant flow.
	http://tools.ietf.org/html/rfc6749#section-4.3
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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

var (
	authReStr = `[Bb]asic (?P<encoded>\S+)`
	authRegex = regexp.MustCompile(authReStr)
)

// PasswordGrant supports authorization via the  Resource Owner Password
// Credentials Grant workflow.
func PasswordGrant(s *Server, w http.ResponseWriter, r *http.Request) {
	//
	// Authenticate
	//
	malformed := "Malformed Authorization header"
	username, password, err := basicAuth(r)
	if err != nil {
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	//
	// Parse authorization request
	//
	dec := json.NewDecoder(r.Body)
	var preq PasswordRequest
	err = dec.Decode(&preq)
	if err != nil && err.Error() != "EOF" {
		msg := "Missing or bad request body"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if username != preq.Username || preq.GrantType != "password" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	t := AuthzTemplate{
		User:   preq.Username,
		Scopes: strings.Split(preq.Scope, " "),
		Note:   preq.Note,
	}
	a, err := s.Authorize(t, password)
	switch {
	case err == ErrNotAuthorized:
		http.Error(w, malformed, http.StatusUnauthorized)
		return
	case err != nil:
		log.Println(err)
		http.Error(w, malformed, http.StatusBadRequest)
	}
	//
	// Authorization granted, compose response
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
