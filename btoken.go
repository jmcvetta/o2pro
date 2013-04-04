// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

// Package btoken implements Oauth2-style bearer tokens.
// https://tools.ietf.org/html/draft-ietf-oauth-v2-bearer-16
package btoken

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"labix.org/v2/mgo"
	"time"
)

const (
	DefaultExpireAfter = "8h" // Duration string for time.ParseDuration()
)

// https://tools.ietf.org/html/draft-ietf-oauth-v2-bearer-16#section-3.1
var (
	// HTTP 400
	ErrInvalidRequest = errors.New("The request is missing a required parameter, includes an unsupported parameter or parameter value, repeats the same parameter, uses more than one method for including an access token, or is otherwise malformed.")
	// HTTP 401
	ErrInvalidToken = errors.New("The access token provided is expired, revoked, malformed, or invalid for other reasons.")
	// HTTP 403
	ErrInsufficientScope = errors.New("The request requires higher privileges than provided by the access token.")
)

// An AuthServer can issue Oauth2-style bearer tokens
type AuthServer interface {
	ExpireAfter(duration string) (time.Duration, error) // Expire all authorizations after duration
	IssueToken(req AuthRequest) (string, error)
	GetAuthorization(token string) (Authorization, error)
}

type AuthRequest struct {
	User     string
	Scopes   []string
	Duration time.Duration
}

type Authorization struct {
	Token      string
	User       string
	Scopes     map[string]bool // Map for easy lookup; always true
	Expiration time.Time
}

// NewMongoAuthServer configures a MongoDB-based AuthServer.  If expireAfter is
// not nil, authorizations will be automatically expired.
func NewMongoAuthServer(db *mgo.Database) (AuthServer, error) {
	m := mongoServer{
		db:   db,
		name: "authorizations",
	}
	_, err := m.ExpireAfter(DefaultExpireAfter)
	return &m, err
}

type mongoServer struct {
	db          *mgo.Database
	name        string // Collection name
	expireAfter time.Duration
}

func (m *mongoServer) ensureIndexes() error {
	//
	// Declare Indexes
	//
	idxs := []mgo.Index{
		mgo.Index{
			Key:      []string{"token"},
			Unique:   true,
			DropDups: false,
		},
		mgo.Index{
			Key:      []string{"expiration"},
			Unique:   true,
			DropDups: false,
		},
	}
	c := m.col()
	for _, i := range idxs {
		err := c.EnsureIndex(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// col returns a Collection object in a new mgo session
func (s *mongoServer) col() *mgo.Collection {
	session := s.db.Session.Copy()
	d := session.DB(s.db.Name)
	return d.C(s.name)
}

func (s *mongoServer) ExpireAfter(duration string) (time.Duration, error) {
	if duration == "" {
		return s.expireAfter, nil
	}
	dur, err := time.ParseDuration("8h")
	if err != nil {
		return dur, err
	}
	s.expireAfter = dur
	err = s.ensureIndexes()
	return dur, err

}

func (s *mongoServer) IssueToken(req AuthRequest) (string, error) {
	c := s.col()
	tok := uuid.NewUUID().String()
	scopes := map[string]bool{}
	dur := req.Duration
	if dur.Seconds() == 0 || dur.Nanoseconds() > s.expireAfter.Nanoseconds() {
		dur = s.expireAfter
	}
	exp := time.Now().Add(dur)
	for _, s := range req.Scopes {
		scopes[s] = true
	}
	a := Authorization{
		Token:      tok,
		User:       req.User,
		Scopes:     scopes,
		Expiration: exp,
	}
	err := c.Insert(a)
	return tok, err
}

func (s *mongoServer) GetAuthorization(token string) (Authorization, error) {
	a := Authorization{}
	c := s.col()
	query := struct {
		Token string
	}{
		Token: token,
	}
	q := c.Find(query)
	cnt, err := q.Count()
	if err != nil {
		return a, err
	}
	if cnt < 1 {
		return a, ErrInvalidToken
	}
	err = q.One(&a)
	if err != nil {
		return a, err
	}
	if time.Now().After(a.Expiration) {
		c.Remove(query)
		return a, ErrInvalidToken
	}
	return a, nil
}

// CheckAuth answers whether the holder of a token is a given user who is
// authorized to access a given scope.  If scope is an empty string, scope is
// not checked.
func (s *mongoServer) CheckAuth(token, user, scope string) (bool, error) {
	a, err := s.GetAuthorization(token)
	if err != nil {
		return false, err
	}
	if a.User != user {
		return false, nil
	}
	if scope != "" {
		_, ok := a.Scopes[scope]
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
