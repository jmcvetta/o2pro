// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"time"
)

const (
	DefaultExpireAfter = "8h" // Duration string for time.ParseDuration()
)

// An AuthServer can issue Oauth2-style bearer tokens
type AuthServer interface {
	ExpireAfter(duration string) (time.Duration, error) // Expire all authorizations after duration
	IssueToken(req AuthRequest) (string, error)
	GetAuthorization(token string) (Authorization, error)
}

