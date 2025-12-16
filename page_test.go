package main


import (
	"testing"
)

func TestReadandWrite(t *testing.T) {

	page := MakePage(10) //page holds 10 bytes (a 10 character string)

	data := []byte("10")

	n, err := page.Write(0, data)

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(data) {
		t.Fatalf("The write returned %d, expected %d", n, len(data))
	}


}