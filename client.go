// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

/*
Client Credentials Grant
https://tools.ietf.org/html/rfc6749#section-4.4
*/

import ()

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

