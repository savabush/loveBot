package main

import (
	"github.com/savabush/breakfastLoveBot/internal/app/breakfastLoveBot"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	go http.ListenAndServe(":6060", nil)
	breakfastLoveBot.App()
}
