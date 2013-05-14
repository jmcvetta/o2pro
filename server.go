// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/rand"
	"encoding/base64"
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

// A Storage back end saves and retrieves authorizations to persistent storage.
type Storage interface {
	saveAuthz(a *Authz) error
	authz(token string) (*Authz, error)
	initialize() error
	migrate() error
}

// An Authenticator authenticates a user's credentials.
type Authenticator func(user, password string) (bool, error)

// A Grantor decides whether to grant access for a given user, scope, and
// client.  Client is optional.
type Grantor func(user, scope string, c *Client) (bool, error)

// GrantAll is a Grantor that always returns true.
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

// SaveAuthz saves an authorization to storage.
func (s *Server) SaveAuthz(a *Authz) error {
	return s.saveAuthz(a)
}

// Authz looks up an authorization based on its token.
func (s *Server) Authz(token string) (*Authz, error) {
	return s.authz(token)
}

// Initialize prepares a fresh database, creating necessary schema, indexes,
// etc.  Behavior is undefined if called with an already-initialized db.
func (s *Server) Initialize() error {
	return s.initialize()
}

// Migrate attempts to update the database to use the latest schema, indexes,
// etc.  Some storage implementations may return ErrNotImplemented.
func (s *Server) Migrate() error {
	return s.migrate()
}

// NewAuth issues a new authorization.
func (s *Server) NewAuthz(t AuthzTemplate) (*Authz, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := uuid.New() + string(b)
	token = base64.StdEncoding.EncodeToString([]byte(token))
	a := Authz{
		Token:      token,
		Uuid:       uuid.New(),
		User:       t.User,
		Scopes:     t.Scopes,
		Expiration: time.Now().Add(s.Duration),
		Note:       t.Note,
	}
	err = s.SaveAuthz(&a)
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
