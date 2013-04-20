// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

// A PasswordRequest is submitted by a client requesting authorization using
// the Resource Owner Password Credentials Grant flow.
// http://tools.ietf.org/html/rfc6749#section-4.3.2
type PasswordRequest struct {
	GrantType string `json:"grant_type"` // REQUIRED.  Value MUST be set to "password".
	Username  string `json:"username"`   // REQUIRED.  The resource owner username.
	Password  string `json:"password"`   // REQUIRED.  The resource owner password.
	Scope     string `json:"scope"`      // OPTIONAL.  The scope of the access request as described by http://tools.ietf.org/html/rfc6749#section-3.3
}
