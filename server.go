// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"code.google.com/p/go-uuid/uuid"
	"log"
	"net/url"
	"os"
	"time"
)

var (
	DefaultExpireAfter = "8h" // Duration string for time.ParseDuration()
	DefaultLogger      = log.New(os.Stdout, "[btoken] ", log.Ltime|log.Ldate|log.Lshortfile)
	DefaultScopes      = []string{"all"}
)

// A Storage back end saves and retrieves authorizations to persistent storage,
// perhaps with caching.
type Storage interface {
	SaveAuth(auth Authorization) error
	GetAuth(token string) (Authorization, error)
	Activate() error // Called when Server is started
}

// A Server is an authorization service that can issue Oauth2-style bearer
// tokens.
type Server struct {
	Storage
	Scopes        []string      // All scopes supported by this server
	DefaultScopes []string      // Issued if no specific scope(s) requested
	MaxDuration   time.Duration // Max lifetime for an authorization
	Logger        *log.Logger
	Authenticator func(url.Userinfo) (bool, error)
}

// NewAuth issues a new Authorization based on an AuthRequest.
func (s *Server) NewAuth(req AuthRequest) (Authorization, error) {
	tok := uuid.NewUUID().String()
	scopes := map[string]bool{}
	dur := req.Duration
	if dur.Seconds() == 0 || dur.Nanoseconds() > s.MaxDuration.Nanoseconds() {
		dur = s.MaxDuration
	}
	exp := time.Now().Add(dur)
	for _, s := range req.Scopes {
		scopes[s] = true
	}
	a := Authorization{
		Token:      tok,
		Owner:      req.Owner,
		Scopes:     scopes,
		Expiration: exp,
	}
	err := s.SaveAuth(a)
	return a, err
}
