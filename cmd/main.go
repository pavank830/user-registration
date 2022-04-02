package main

import (
	"flag"
	"fmt"
	"log"

	userregistration "github.com/pavank830/user-registration"
)

func main() {
	var serverPortFlag = flag.String("port", "8080", "HTTP server port,default port 8080")
	var dbAdrrFlag = flag.String("db", "", "mysql database connection in <username>:<pwd>@tcp(<host>:<port>)/<db>")
	flag.Parse()
	fmt.Println("HTTP server port value:", *serverPortFlag)
	fmt.Println("sql database connection:", *dbAdrrFlag)
	if *dbAdrrFlag == "" {
		log.Fatalln("sql database connection addr is empty")
	}
	userregistration.Start(*serverPortFlag, *dbAdrrFlag)
}
