package main

import (
	"bird-watcher/internal/handlers"
	watcher "bird-watcher/internal/services"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func HandleCli() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit", text) == 0 {
			log.Println("bird-watcher going offline.")
			os.Exit(0)
		}
	}
}

func main() {
	log.Println("bird-watcher online and scanning at port 8080.")

	// handle cli inputs
	go HandleCli()

	// start watcher job
	watcher.StartWatcher()

	// manage handlers
	http.HandleFunc("/", handlers.Index)

	// get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// start server
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("failed to boot up server.")
	}
}
