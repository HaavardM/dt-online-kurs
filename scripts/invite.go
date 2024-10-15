package main

import (
	"dt-online-kurs/internal/disruptive"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: invite <email>")
	}

	err := disruptive.InviteUser(os.Args[1])
	if err != nil {
		panic(err)
	}
}
