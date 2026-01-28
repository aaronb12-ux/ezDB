package main

import (
	"log"
	"fmt"
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

	val, _ := db.Get(1)

	fmt.Printf("we got the value %v", string(val))

	e2 := db.Put(1, []byte("cosmo"))

	if e2 != nil {
		log.Fatal("error", e2)
	}

}





