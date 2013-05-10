// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Authorization struct {
	Id         bson.ObjectId `bson:"_id",json:"id"` // Unique storage-dependent ID for this Authorization
	Token      string        `json:"token"`
	Username   string        `json:"username"`
	Scopes     []string      `json:"scopes"`
	Expiration time.Time     `json:"expiration"`
	Note       string        `json:"note"`
}

type ClientType string

const (
	PublicClient       ClientType = "public"
	ConfidentialClient            = "Confidential"
)

type Client struct {
	Id          bson.ObjectId `bson:"_id",json:"id"` // Storage-dependent ID for this Client
	ClientType  ClientType
	RedirectUri string
	AppName     string
	WebSite     string
	Description string
}
