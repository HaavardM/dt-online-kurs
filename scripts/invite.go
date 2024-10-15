package main

import (
	"os"

	"github.com/HaavardM/dt-online-kurs/internal/disruptive"
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
