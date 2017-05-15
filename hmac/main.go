package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

const usage = `
usage:
	hmac sign|verify <key> <value>
`

func main() {
	if len(os.Args) < 4 ||
		(os.Args[1] != "sign" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1] //sign or verify
	key := os.Args[2]
	value := os.Args[3]

	switch cmd {
	case "sign":
		v := []byte(value)
		h := hmac.New(sha256.New, []byte(key)) //create a new hash authentication code
		h.Write(v)
		sig := h.Sum(nil) //takes all the data written to the hashing algorithm and computes digsig
		// fmt.Println(sig)

		//Combine the original data + the digsig -> base64 encode and send across network
		buf := make([]byte, len(v)+len(sig)) //How big the buffer is
		copy(buf, v)                         //Take v bytes and copy into buf
		copy(buf[len(v):], sig)              //Take a slice of the buf slice
		fmt.Println(base64.URLEncoding.EncodeToString(buf))
	case "verify":
		buf, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			os.Exit(1)
		}
		//We know how long the hash will be because sha56
		v := buf[:len(buf)-sha256.Size]
		sig := buf[len(buf)-sha256.Size:]

		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		sig2 := h.Sum(nil)
		if hmac.Equal(sig, sig2) { //.Equal negates timing attacks
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Boooo invalid signature")
		}
	}
}
