package filemanager


import (
	"testing"
	"fmt"
)

func TestWriteWhole(t *testing.T) {

	page := MakePage(10) //page holds 10 bytes (a 10 character string)

	data := []byte("hiii")

	n, err := page.Write(0, data)
	
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(data) {
		t.Fatalf("The write returned %d, expected %d", n, len(data))
	}
}

func TestReadWhole(t *testing.T) {

	page := MakePage(5) //first write to a page

	n , err := page.Write(0, []byte("Cosmo"))

	if err != nil {
		t.Fatalf("Write to page failed: %v", err)
	}

	if n != 5 {
		t.Fatalf("The write returned %d, expected 5", n)
	}

	//now read
	data := make([]byte, 5)

	h := page.Read(0, data) //read data from page and put it into 'data'

	if h != 5 {
		t.Fatalf("The write returned %d but we expected 5", h)
	}
}

func TestWriteOffSet(t *testing.T) {

	page := MakePage(20)

	data := []byte("Cosmo, is awesome!")

	n, err := page.Write(2, data) //writing to page starting at offset 4. When we offset, we write from [offset : end of page].

	if err != nil {
		t.Fatalf("Write to page failed %v", err)
	}

	fmt.Printf("We wrote %d bytes to the page", n)

	if n != len(data) {
		t.Fatalf("The write returned %d but we expected %d", n, len(data))
	}
}

func TestReadOffSet(t *testing.T) {

	//first write to the page with an offset

	page := MakePage(10) //page stores 10 bytes

	data := []byte("Cosmo")

	n, err := page.Write(3, data) //writing to page starting at offset 3 

	if err != nil {
		t.Fatalf("Error writing to page %v", err)
	}

	if n != len(data) {
		t.Fatalf("The write returned %d, but we expected %d", n, len(data))
	}

	//now read from the page at offset and store it

	store := make([]byte, 5) //stores 5 bytes

	k := page.Read(3, store)

	if k != 5 { //we are expected to read 5 bytes in page, because it has 5 bytes
		t.Fatalf("The write returned %d, but we expected %d", k, 5)
	}

	fmt.Printf("'store' is: %v", string(store))
}