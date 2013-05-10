// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

/*
http://tools.ietf.org/html/rfc6749#section-2
*/

import (
	"labix.org/v2/mgo/bson"
)

type ClientType int

const (
	PublicClient ClientType = iota
	ConfidentialClient
)

type ClientTemplate struct {
}

type Client struct {
	Id          bson.ObjectId `bson:"_id",json:"id"` // Storage-dependent ID for this Client
	ClientType  ClientType
	RedirectUri string
	AppName     string
	WebSite     string
	Description string
}
