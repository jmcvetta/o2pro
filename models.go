// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// An Authz is an authorization.
type Authz struct {
	Id         int64
	Uuid       string `bson:"_id"`
	Token      string
	Username   string
	Scopes     []string
	Issued     time.Time
	Expiration time.Time
	Note       string
}

type ClientType string

const (
	PublicClient       ClientType = "public"
	ConfidentialClient            = "confidential"
)

type Client struct {
	Id          bson.ObjectId `bson:"_id",json:"id"` // Storage-dependent ID for this Client
	ClientType  ClientType
	RedirectUri string
	AppName     string
	WebSite     string
	Description string
}
