// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package btoken

import (
	"errors"
)

// https://tools.ietf.org/html/draft-ietf-oauth-v2-bearer-16#section-3.1
var (
	// HTTP 400
	ErrInvalidRequest = errors.New("The request is missing a required parameter, includes an unsupported parameter or parameter value, repeats the same parameter, uses more than one method for including an access token, or is otherwise malformed.")
	// HTTP 401
	ErrInvalidToken = errors.New("The access token provided is expired, revoked, malformed, or invalid for other reasons.")
	// HTTP 403
	ErrInsufficientScope = errors.New("The request requires higher privileges than provided by the access token.")
)
