// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package main

import (
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/btoken"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
)

var (
	addr = ":8080"
)

// fakeAuth authorizes everyone for everything.
func fakeAuth(*url.Userinfo, btoken.AuthRequest) (bool, error) {
	return true, nil
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// Connect to MongoDB
	//
	mongoUrl := env.StringDefault("MONGOLAB_URI", "localhost")
	log.Println("Connecting to MongoDB on " + mongoUrl + "...")
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	db := session.DB("")
	_, err = db.CollectionNames()
	if err != nil || db.Name == "test" {
		log.Println("Setting db name to 'btoken'.")
		db = session.DB("btoken")
	}
	a := btoken.Authorizer(fakeAuth)
	srv, err := btoken.NewMongoServer(db, "8h", a)
	if err != nil {
		log.Fatal(err)
	}
	hf := srv.AuthReqHandler()
	http.HandleFunc("/auth", hf)
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
