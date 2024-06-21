package main

import (
	"fmt"

	"github.com/stack1ng/azuretls-client"
)

func main() {
	session := azuretls.NewSession()

	session.OrderedHeaders = azuretls.OrderedHeaders{
		{"user-agent", "test"},
		{"content-type", "application/json"},
		{"accept", "application/json"},
	}

	response, err := session.Get("https://tls.peet.ws/api/all")

	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))

	session.Close()
}
