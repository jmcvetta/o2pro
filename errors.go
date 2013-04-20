// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

import (
	"errors"
)

// Standard Oauth2 error types
// https://tools.ietf.org/html/draft-ietf-oauth-v2-bearer-16#section-3.1
var (
	// HTTP 400
	ErrInvalidRequest = errors.New("The request is missing a required parameter, includes an unsupported parameter or parameter value, repeats the same parameter, uses more than one method for including an access token, or is otherwise malformed.")
	// HTTP 401
	ErrNotAuthorized = errors.New("Authorization not granted.")
	ErrInvalidToken  = errors.New("The access token provided is expired, revoked, malformed, or invalid for other reasons.")
	// HTTP 403
	ErrInsufficientScope = errors.New("The request requires higher privileges than provided by the access token.")
)

type oauthError string

// Error Codes
const (
	invalidRequest = oauthError("invalid_request")
	/*
		The request is missing a required parameter, includes an
		unsupported parameter value (other than grant type),
		repeats a parameter, includes multiple credentials,
		utilizes more than one mechanism for authenticating the
		client, or is otherwise malformed.
	*/

	invalidClient = "invalid_client"
	/*
		Client authentication failed (e.g., unknown client, no
		client authentication included, or unsupported
		authentication method).  The authorization server MAY
		return an HTTP 401 (Unauthorized) status code to indicate
		which HTTP authentication schemes are supported.  If the
		client attempted to authenticate via the "Authorization"
		request header field, the authorization server MUST
		respond with an HTTP 401 (Unauthorized) status code and
		include the "WWW-Authenticate" response header field
		matching the authentication scheme used by the client.
	*/

	invalidGrant = "invalid_grant"
	/*
		The provided authorization grant (e.g., authorization
		code, resource owner credentials) or refresh token is
		invalid, expired, revoked, does not match the redirection
		URI used in the authorization request, or was issued to
		another client.
	*/

	unauthorizedClient = "unauthorized_client"
	/*
		The authenticated client is not authorized to use this
		authorization grant type.
	*/

	unsupportedGrantType = "unsupported_grant_type"
	/*
		The authorization grant type is not supported by the
		authorization server.
	*/

	invalidScope = "invalid_scope"
	/*
		The requested scope is invalid, unknown, malformed, or
		exceeds the scope granted by the resource owner.
	*/
)

// An ErrorResponse is sent with HTTP status code 400.
type ErrorResponse struct {
	Error string `json:"error"`                       // REQUIRED.  A single ASCII error code from the Error Codes constants.
	Desc  string `json:"error_description,omitempty"` // OPTIONAL.  Human-readable ASCII [USASCII] text providing additional information, used to assist the client developer in understanding the error that occurred.
	Uri   string `json:"error_uri,omitempty"`         // OPTIONAL.  A URI identifying a human-readable web page with information about the error, used to provide the client developer with additional information about the error.

}
