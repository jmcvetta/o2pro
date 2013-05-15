// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package o2pro

/*
Implementation of AUTHORIZATION CODE GRANT workflow
https://tools.ietf.org/html/rfc6749#section-4.1
 */

import (
)

// A Code is an authorization code, entitling its holder to be issued an
// authorization.
type Code struct {
	Id int64 `bson:",omitempty`
}
