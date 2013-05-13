// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	DefaultExpireAfter = "8h" // Duration string for time.ParseDuration()
	DefaultLogger      = log.New(os.Stdout, "[o2pro] ", log.Ltime|log.Ldate|log.Lshortfile)
	DefaultScopes      = []string{"all"}
)

type Scopes map[string]bool // scope_name:true

// A Storage back end saves and retrieves authorizations to persistent storage.
type Storage interface {
	SaveAuthz(a *Authz) error
	Authz(token string) (*Authz, error)
	// Run() error // Called when Server is started
	Initialize() error
	Migrate() error
}

// An Authenticator authenticates a user's credentials.
type Authenticator func(user, password string) (bool, error)

// A Grantor decides whether to grant access for a given user, scope, and
// client.  Client is optional.
type Grantor func(user, scope string, c *Client) (bool, error)

// GrantAll returns true regardless of scope and client.
func GrantAll(user, scope string, c *Client) (bool, error) {
	return true, nil
}

func NewServer(s Storage, a Authenticator, g Grantor) *Server {
	dur, err := time.ParseDuration(DefaultExpireAfter)
	if err != nil {
		log.Panic(err)
	}
	return &Server{
		Storage:       s,
		Scopes:        DefaultScopes,
		DefaultScopes: DefaultScopes,
		Duration:      dur,
		Logger:        DefaultLogger,
		a:             a,
		g:             g,
	}
}

// A Server is an OAuth2 authorization server.
type Server struct {
	Storage
	Scopes        []string      // All scopes supported by this server
	DefaultScopes []string      // Issued if no specific scope(s) requested
	Duration      time.Duration // Lifetime for an authorization
	Logger        *log.Logger
	a             Authenticator
	g             Grantor
}

// Grant decides whether to grant an authorization.
func (s *Server) Grant(user, scope string, c *Client) (bool, error) {
	return s.g(user, scope, c)
}

// Authenticate validates a user's credentials.
func (s *Server) Authenticate(user, password string) (bool, error) {
	return s.a(user, password)
}

// NewAuth issues a new authorization.
func (s *Server) NewAuthz(t AuthzTemplate) (*Authz, error) {
	a := Authz{
		Token:      uuid.New(),
		Uuid:       uuid.New(),
		User:       t.User,
		Scopes:     t.Scopes,
		Expiration: time.Now().Add(s.Duration),
		Note:       t.Note,
	}
	err := s.SaveAuthz(&a)
	return &a, err
}

type handlerStub func(s *Server, w http.ResponseWriter, r *http.Request)

func (s *Server) HandlerFunc(hs handlerStub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hs(s, w, r)
	}
}

/*
func (s *Server) Error(w http.ResponseWriter, error string, code int) {

}
*/

/*
// Authorize may grant an authorization to a client.  Server.Authorizer
// decides whether to make the grant. ErrNotAuthorized is returned if
// authorization is denied.
func (s *Server) Authorize(t AuthzTemplate, password string) (*Authz, error) {
	a := new(Authz)
	ok, err := s.Authorizer(t.User, password, t.Scopes)
	if err != nil {
		return a, err
	}
	if !ok {
		return a, ErrNotAuthorized
	}
	return s.NewAuthz(t)
}
*/
