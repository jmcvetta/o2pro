// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released
// under the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for
// details.  Resist intellectual serfdom - the ownership of ideas is akin to
// slavery.

package btoken

import (
	"time"
)

type AuthRequest struct {
	User     string
	Scopes   []string
	Duration time.Duration // Max duration is AuthServer.ExpireAfter
}

type Authorization struct {
	Token      string
	User       string
	Scopes     map[string]bool // Map for easy lookup; always true
	Expiration time.Time
}
