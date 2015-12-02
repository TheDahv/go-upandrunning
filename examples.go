package main

// This little program determine which example the user wants to run, loads it
// up, and feeds it arguments. In the case of more complex examples, like the
// New York Times API demonstrations, we also perform error handling, results
// reading, and printing to the console.

// By convention and style, imports are grouped into one statement and lexically
// sorted. Standard library packages are grouped first, followed by packages
// from the same library, and then 3rd-party packages.
// `goimports` is managing all this for us!
import (
	"bufio"
	"fmt"
	"os"

	"github.com/thedahv/go-upandrunning/examples/hello"
	"github.com/thedahv/go-upandrunning/examples/nyt"
)

func main() {
	// A simple example of reading arguments from the program. Note, these are
	// actual arguments, and not stdin (that is accessible a different way as an
	// io.Reader)
	//
	// Since the first argument is our program itself, we can "take a slice" of
	// our arguments to skip the first entry.
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: ./examples [example] [args...]")
		fmt.Println("Available examples are: hello, nyt-single, nyt-serial, nyt-parallel")

		return
	}

	switch args[0] {
	case "hello":
		hello.Run()
	case "nyt-single", "nyt-serial", "nyt-parallel":
		if len(args) > 1 {
			runNYTExample(args)
		} else {
			fmt.Println("Need a search term for NYT example")
		}
	default:
		fmt.Printf("No recognizable examples for '%s'\n", args[0])
	}
}

// A little helper to run our NYT examples, parse output, and write to Stdout
func runNYTExample(args []string) {

	// All the NYT examples have basically the same signature. So we lifted the
	// signature into a type called "Runnable", and declare it as a variable. Go
	// functions are first-class values, and can be assigned and passed around.
	var runner nyt.Runnable

	switch args[0] {
	case "nyt-single":
		runner = nyt.RunSingle
	case "nyt-serial":
		runner = nyt.RunSerial
	case "nyt-parallel":
		runner = nyt.RunParallel
	}

	results, err := runner(args[1:])

	if err != nil {
		fmt.Println(err)
		return
	}

	if results == nil {
		fmt.Println("No results for that query")
		return
	}

	// Scanner is a great way to consume an io.Reader with buffered reads. Note,
	// it splits a reader by newlines by default, but you can create a scanner
	// that splits on whatever you want.
	scanner := bufio.NewScanner(results)
	for scanner.Scan() {
		// Nothing too interesting here. Just sending results back to Stdout. We
		// could have done the same with fmt.Println.
		// It's worth noting bufio.Scanner eats the delimiting character, so we add
		// our newline back in below.
		fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
	}
}
