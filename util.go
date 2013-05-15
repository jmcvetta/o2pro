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

func BasicAuth(r *http.Request) (username, password string, err error) {
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
func BearerToken(r *http.Request) (token string, err error) {
	//
	// Authorization Header
	//
	str := r.Header.Get("Authorization")
	if str != "" {
		matches := bearerRegex.FindStringSubmatch(str)
		if len(matches) != 2 {
			log.Println("Regex doesn't match")
			err = ErrInvalidRequest
			return
		}
		token = matches[1]
		return
	}
	//
	// Form-encoded Body Parameter
	//
	return "", ErrNotImplemented
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
