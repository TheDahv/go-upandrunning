# Go: Up and Running

This is a quick overview of the Go programming language. It aims to be a quick tour
of language features and simple guide to get you to a bare-bones programming
environment.

It makes use of the [presentation
tool](https://godoc.org/golang.org/x/tools/present) baked into Go. That means
you can run this presentation on your own laptop if you already have Go
installed.

If you do, clone this repository and run the following from within it:

* `go get golang.org/x/tools/cmd/present`
* `present`

It should tell you what port to open to view the presentation.

## Topics Covered

* Basic Hello World
* Simple language features (structs, loops, errors, slices)
* HTTP API calls and JSON processing
* Parallelizing a series of synchronous IO calls with goroutines and channels

## Possible Future Topics

* Using baked-in profiling tools to analyze performance
* Organizing code with structs and interfaces

# Examples

The slides include code lifted from examples you can run yourself. You can
run `go run examples.go` to see what examples are available and how to run them.

If you want to run the New York Times API examples, you'll need to register
for an API Key from the
[New York Times Developers](http://developer.nytimes.com/) site.

Either add your key to the `API_KEY` environment variable in your session, or
add it in a `.key` file at the root of this directory and `./run.sh` will
inject it into your program.

# Get Involved!

Did you find something fishy? Do you see something that should be covered that
isn't?

Feel free to add something in the Issues and I'd be happy to take a look!
