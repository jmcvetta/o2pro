// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"time"
)

// An AuthRequest describes the details of an Authorization to be issued.
type AuthRequest struct {
	Scopes   []string      // http://tools.ietf.org/html/rfc6749#section-3.3
	Duration time.Duration // Max duration is AuthServer.ExpireAfter
	Note     string        // Optional
}

type Authorization struct {
	AuthId     string // Unique storage-dependent ID for this Authorization
	Token      string
	Owner      string
	Scopes     map[string]bool // Map for easy lookup; always true
	Expiration time.Time
	Note       string
}
