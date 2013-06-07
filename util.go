// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	basicRegex  = regexp.MustCompile(`[Bb]asic (?P<encoded>\S+)`)
	bearerRegex = regexp.MustCompile(`[Bb]earer (?P<token>\S+)`) // Spec doesn't actually say "Bearer" should be case insensitive.
)

// BasicAuth extracts username & password from an HTTP request's authorization
// header.
func basicAuth(r *http.Request) (username, password string, err error) {
	str := r.Header.Get("Authorization")
	matches := basicRegex.FindStringSubmatch(str)
	if len(matches) != 2 {
		log.Println("Regex doesn't match")
		err = ErrInvalidRequest
		return
	}
	encoded := matches[1]
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		log.Println("Base64 decode failed")
		err = ErrInvalidRequest
		return
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		log.Println("String split failed")
		err = ErrInvalidRequest
		return
	}
	username = parts[0]
	password = parts[1]
	return
}

// BearerToken extracts a bearer token from the authorization header, form
// encoded body parameter, or URI query parameter of an HTTP request.
func bearerToken(r *http.Request) (token string, err error) {
	//
	// Authorization Header
	//
	auth := r.Header.Get("Authorization")
	if auth != "" {
		matches := bearerRegex.FindStringSubmatch(auth)
		if len(matches) != 2 {
			log.Println("Regex doesn't match")
			log.Println("\t" + auth)
			err = ErrNoToken
			return
		}
		token = matches[1]
		return
	}
	//
	// Form-encoded Body Parameter
	//
	ct := r.Header.Get("Content-Type")
	if ct == "application/x-www-form-urlencoded" {
		type t struct {
			Token string `json:"access_token"`
		}
		s := new(t)
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err = dec.Decode(s)
		token = s.Token
		return
	}
	token = r.URL.Query().Get("access_token")
	if token == "" {
		err = ErrNoToken
	}
	return
}

func sliceMap(s []string) map[string]bool {
	sm := make(map[string]bool, len(s))
	for _, s := range s {
		sm[s] = true
	}
	return sm
}

func prettyPrint(v interface{}) {
	_, file, line, _ := runtime.Caller(1)
	lineNo := strconv.Itoa(line)
	file = filepath.Base(file)
	b, _ := json.MarshalIndent(v, "", "\t")
	s := file + ":" + lineNo + ": \n" + string(b) + "\n"
	println(s)
}
