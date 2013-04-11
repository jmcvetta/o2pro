// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"encoding/json"
	"log"
)

var (
	authReStr = `[Bb]asic (?P<encoded>\D+)`
	authRegex = regexp.MustCompile(authReStr)
)

func (s *Server) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	//
	// Authenticate
	//
	str := r.Header.Get("Authorization")
	malformed := "Malformed Authorization header"
	matches := authRegex.FindStringSubmatch(str)
	if len(matches) != 2 {
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	encoded := matches[1]
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	username := parts[0]
	password := parts[1]
	u := url.UserPassword(username, password)
	//
	// Parse authorization request
	//
	dec := json.NewDecoder(r.Body)
	var areq AuthRequest
	err = dec.Decode(&areq)
	if err != nil {
		log.Println(err)
		http.Error(w, malformed, http.StatusBadRequest)
		return
	}
	a, err := s.Authorize(u, areq)
	switch {
		case err == ErrUnauthorized:
			http.Error(w, malformed, http.StatusUnauthorized)
			return
		case err != nil:
			log.Println(err)
			http.Error(w, malformed, http.StatusBadRequest)
	}
	//
	// Authorization Granted
	//
	enc := json.NewEncoder(w)
	enc.Encode(&a)
	return
}
