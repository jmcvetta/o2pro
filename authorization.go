// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"labix.org/v2/mgo/bson"
	"time"
)

/*
TODO: Make token req/resp conform to RFC:
	http://tools.ietf.org/html/rfc6749#section-4.3.2
	http://tools.ietf.org/html/rfc6749#section-4.3.3
*/

// An AuthRequest describes the details of an Authorization to be issued.
type AuthRequest struct {
	Scopes      []string      // http://tools.ietf.org/html/rfc6749#section-3.3
	ExpireAfter time.Duration // Max duration is AuthServer.ExpireAfter
	Note        string        // Optional
}

// http://tools.ietf.org/html/rfc6749#section-4.3.2
type TokenRequest struct {
	GrantType string `json:"grant_type"` // REQUIRED.  Value MUST be set to "password".
	Username  string `json:"username"`   // REQUIRED.  The resource owner username.
	Password  string `json:"password"`   // REQUIRED.  The resource owner password.
	Scope     string `json:"scope"`      // OPTIONAL.  The scope of the access request as described by http://tools.ietf.org/html/rfc6749#section-3.3
}

type Authorization struct {
	AuthId     bson.ObjectId   `bson:"_id",json:"id"` // Unique storage-dependent ID for this Authorization
	Token      string          `json:"token"`
	Owner      string          `json:"owner"`
	Scopes     []string        `json:"scopes"`
	ScopesMap  map[string]bool `json:"-"` // Map for easy lookup; always true
	Expiration time.Time       `json:"expiration"`
	Note       string          `json:"note"`
}
