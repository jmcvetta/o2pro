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
	"runtime"
	"strconv"
	"strings"
)

func basicAuth(r *http.Request) (username, password string, err error) {
	str := r.Header.Get("Authorization")
	matches := authRegex.FindStringSubmatch(str)
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

func prettyPrint(v interface{}) {
	_, file, line, _ := runtime.Caller(1)
	lineNo := strconv.Itoa(line)
	file = filepath.Base(file)
	b, _ := json.MarshalIndent(v, "", "\t")
	s := file + ":" + lineNo + ": \n" + string(b) + "\n"
	println(s)
}
