// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"time"
)

// An Authz is an authorization.
type Authz struct {
	Id         int64 `bson:",omitempty`
	Uuid       string
	Token      string
	User       string
	ClientId   int64
	Client     *Client
	Issued     time.Time
	Expiration time.Time
	Note       string
	Scopes     []string
}

/*
type AuthzScope struct {
	Id      int64
	AuthzId int64 `qbs:"fk:Authz`
	Authz   *Authz
	Scope   string
}
*/

const (
	PublicClient       = "public"
	ConfidentialClient = "confidential"
)

// A Client is an application making protected resource requests on behalf of
// the resource owner and with its authorization.
type Client struct {
	Id          int64  `bson:",omitempty`
	ClientType  string // "public" or "confidential"
	RedirectUri string
	AppName     string
	WebSite     string
	Description string
}

// A Code is an authorization code, entitling its holder to be issued an
// authorization.
type Code struct {
	Id int64 `bson:",omitempty`
}
