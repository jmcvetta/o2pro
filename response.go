// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

// A TokenResponse is sent on a successful authorization request.
// https://tools.ietf.org/html/rfc6749#section-5.1
type TokenReponse struct {
	AccessToken  string `json:"access_token"`            // REQUIRED.  The access token issued by the authorization server.
	TokenType    string `json:"token_type"`              // REQUIRED.  The type of the token issued as described in Section 7.1.  Value is case insensitive.
	ExpiresIn    int    `json:"expires_in,omitempty"`    // RECOMMENDED.  The lifetime in seconds of the access token.  For example, the value "3600" denotes that the access token will expire in one hour from the time the response was generated. If omitted, the authorization server SHOULD provide the expiration time via other means or document the default value.
	RefreshToken string `json:"refresh_token,omitempty"` //  OPTIONAL.  The refresh token, which can be used to obtain new access tokens using the same authorization grant as described in Section 6.
	Scope        string `json:"scope"`                   //  OPTIONAL, if identical to the scope requested by the client; otherwise, REQUIRED.  The scope of the access token as described by Section 3.3.
}
