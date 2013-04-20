// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package main

import (
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/o2pro"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
)

var (
	addr = ":8080"
)

// fakeAuth authorizes everyone for everything.
func fakeAuth(*url.Userinfo, o2pro.AuthRequest) (bool, error) {
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
		log.Println("Setting db name to 'o2pro'.")
		db = session.DB("o2pro")
	}
	a := o2pro.Authorizer(fakeAuth)
	srv, err := o2pro.NewMongoServer(db, "8h", a)
	if err != nil {
		log.Fatal(err)
	}
	hf := srv.AuthReqHandler()
	http.HandleFunc("/auth", hf)
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
