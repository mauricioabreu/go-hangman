# Hangman

[![CircleCI](https://circleci.com/gh/mauricioabreu/go-hangman.svg?style=svg)](https://circleci.com/gh/mauricioabreu/go-hangman)

Hangman game written in golang.

`Hangman` is a very famous game that a player must guess
the word by suggesting letters being limited by a low number of
guesses.

More information here: https://en.wikipedia.org/wiki/Hangman_(game)

## Architecture

The game is in development phase but the architecture we are aiming looks like the following:

![Hangman architecture overview, with all components connected](misc/hangman_architecture.png "Architecture overview")

An API provides an elegant way to query and mutate the game state. Any client (command line interfaces, browsers) can connect to this API which talks to the backend server, executing the actions needed. There is also a storage that persists every move.

## Installing

To install the project you can run the following command:

```
go get github.com/mauricioabreu/go-hangman
```

To run the test suite, use the command `go test -v ./...`

More information is available in our [contributing](CONTRIBUTING.md) guide.

### Dependencies

You can use the [dep](https://golang.github.io/dep/) tool to build dependencies. They are alredy vendored in the project, but
dep will make it easier to add new dependencies and pin their versions.

Follow the installation guide. Then you can use the command `dep ensure` to install the dependencies.

It may take a while. After that, all the dependencies will be installed inside `vendor`.

Use the command `dep status` to check every dependency and its version.

## ReST API

A ReST API is being developed to provide a consistent way to develop a new client based on the server.

To interact with the server, you first need to run it:

    go run http/api.go

After having the server up and running you can start to talk to the API.

To specify a words file you want to load, you can use the `-words_file` command line option:

    go run http/api.go -words_file=/path/to/my/file.txt

There is an in-memory storage for this game, so you won't need to run a real database:

    go run http/api.go -datastore memory

Other flags can be listed using the help flag:

    go run http/api.go -h

### Starting a new game:

    curl -sv http://localhost:8000/games -XPOST

This endpoit returns a header `Location` with a new endpoint, used to make guesses.

### Checking the current status of a game:

The `ID` in the `Location` header can be used to know the current state of the game:

    curl -sv http://localhost:8000/games/e46964a6-cd18-4820-b208-449ddbeb2d83 -XGET

A JSON response will be returned, containing a pretty representation of the game. If the game does not exists, HTTP 404 status will be returned otherwise.

### Guessing a letter for a word:

You can now guess a letter for the word by sending a JSON payload to the server:

    curl -sv http://localhost:8000/games/e46964a6-cd18-4820-b208-449ddbeb2d83/guesses -XPUT -H "Content-Type: application/json" -d '{"guess": "a"}'

### Removing a game:

If a game need to be deleted for some reason, there is a way to do that:

    curl -sv http://localhost:8000/games/e46964a6-cd18-4820-b208-449ddbeb2d83 -XDELETE -H "Content-Type: application/json"

If the game exists, HTTP 204 is returned - if the game does not exist HTTP 404 is returned. 
