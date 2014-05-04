package main

import (
	"fmt"
	"os"
)

func ReadAndEncrypt(filename string) {
	// read a file in chunks, encrypt, send!
	// Errors here should maybe result in pushing those event back onto the queue?
	// or push them onto an error channel so that they can be handled else where?

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err, "\n\n")
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	fmt.Println(stat.Size())

	var amount int64 = 0
	var EOF bool = false
	for !EOF {
		if amount >= stat.Size() {
			EOF = true
		}

		data := make([]byte, 16)
		count, err := file.Read(data)

		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		amount += int64(count)
		fmt.Printf("read %d bytes: %q\n", count, data[:count])
	}



}
