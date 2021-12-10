package main

import "github.com/fernandoocampo/users-micro/internal/application"

func main() {
	newInstance := application.NewInstance()
	err := newInstance.Run()
	if err != nil {
		panic(err)
	}
	defer newInstance.Stop()
}
