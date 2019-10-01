/*
Use html/template to create HTML Pages
Create an http.Handler to handle the web requests instead of a handler function
Use encoding/json to decode the JSON File

Story starts at intro

1. Decode JSON File
2. Store JSON File into a map
3. Create HTML template
4. Create handler
5. Serve files
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yousseffarkhani/gophercises/03-chooseYourOwnAdventure/cyoa"
)

func main() {
	port := flag.Int("port", 8080, "Port to start the web application")
	filename := flag.String("file", "gopher.json", "JSON File for the CYOA story")
	flag.Parse()

	file, err := os.Open(*filename)
	checkError(err)
	defer file.Close()

	story, err := cyoa.JsonStory(file)
	checkError(err)

	handler := cyoa.NewHandler(story)                                  // Avec les functional options, on peut passer un nouveau template avec WithTemplate
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler)) // Ne retourne rien si tout va bien. Retourne une erreur en cas de problème.
	fmt.Printf("Starting the server on port: %d\n", *port)

}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
