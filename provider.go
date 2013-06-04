// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
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

// NewProvider initializes a new OAuth2 provider server.
func NewProvider(s Storage, a Authenticator, g Grantor) *Provider {
	dur, err := time.ParseDuration(DefaultExpireAfter)
	if err != nil {
		log.Panic(err)
	}
	return &Provider{
		Storage:       s,
		Scopes:        DefaultScopes,
		DefaultScopes: DefaultScopes,
		Duration:      dur,
		Logger:        DefaultLogger,
		a:             a,
		g:             g,
	}
}

// A Provider is an OAuth2 authorization server.
type Provider struct {
	Storage
	Scopes        []string      // All scopes supported by this server
	DefaultScopes []string      // Issued if no specific scope(s) requested
	Duration      time.Duration // Lifetime for an authorization
	Logger        *log.Logger
	a             Authenticator
	g             Grantor
}

// Grant decides whether to grant an authorization.
func (p *Provider) Grant(user, scope string, c *Client) (bool, error) {
	return p.g(user, scope, c)
}

// Authenticate validates a user's credentials.
func (p *Provider) Authenticate(user, password string) (bool, error) {
	return p.a(user, password)
}

// Initialize prepares a fresh database, creating necessary schema, indexes,
// etc.  Behavior is undefined if called with an already-initialized db.
func (p *Provider) Initialize() error {
	return p.initialize()
}

// Migrate attempts to update the database to use the latest schema, indexes,
// etc.  Some storage implementations may return ErrNotImplemented.
func (p *Provider) Migrate() error {
	return p.migrate()
}

type handlerStub func(p *Provider, w http.ResponseWriter, r *http.Request)

func (p *Provider) HandlerFunc(hs handlerStub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hs(p, w, r)
	}
}
