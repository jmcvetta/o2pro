// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"time"
)

// An Authz is an authorization.
type Authz struct {
	Id         int64  `bson:",omitempty`
	Uuid       string `bson:"_id"`
	Token      string
	User       string // Unique user identifier
	Client     *Client
	ClientId   int64 `bson:",omitempty"`
	Scopes     []string
	Issued     time.Time
	Expiration time.Time
	Note       string
}

type Scope struct {
	Id int64 `bson:",omitempty`
}

type ClientType string

const (
	PublicClient       ClientType = "public"
	ConfidentialClient            = "confidential"
)

type Client struct {
	Id          int64 `bson:",omitempty`
	ClientType  ClientType
	RedirectUri string
	AppName     string
	WebSite     string
	Description string
}
