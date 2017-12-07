package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/SpaceHexagon/ecs/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("| Welcome to ECS, %s\n",
		user.Username)
	fmt.Printf("| Interactive Mode\n")
	repl.Start(os.Stdin, os.Stdout)
}
