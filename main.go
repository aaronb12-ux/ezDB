package main

import (
	"log"
	
)

func main() {

	db, err := Open("hi.db", 4096, 4)

	if err != nil {
		log.Fatalf("error creating db")
	}

	e1 := db.Put(1, []byte("coolness"))

	if e1 != nil {
		log.Fatal("error", e1)
	}

	e2 := db.Put(2, []byte("cosmo"))

	if e2 != nil {
		log.Fatal("error", e2)
	}

	e3 := db.Put(3, []byte("lol"))

	if e3 != nil {
		log.Fatal("error", e3)
	}

	ans = db.Get(3)

	db.ShowAll()

}





