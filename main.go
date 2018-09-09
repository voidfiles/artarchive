package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/voidfiles/artarchive/cmd"
)

func main() {
	if lambdaIs := os.Getenv("LAMBDA"); lambdaIs != "" {
		log.Printf("Running lambda: %s", lambdaIs)
		lambda.Start(func() (string, error) {
			cmd.Execute(lambdaIs)
			return fmt.Sprintf("We made it"), nil
		})
	} else {
		cmd.Execute("")
	}

}
