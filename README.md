# Hangman

Hangman game written in golang.

`Hangman` is a very famous game that a player must guess
the word by suggesting letters being limited by a low number of
guesses.

More information here: https://en.wikipedia.org/wiki/Hangman_(game)

## Architecture

The game is in development phase but the architecture we are aiming looks like the following:

![Hangman architecture overview, with all components connected](misc/hangman_architecture.png "Architecture overview")

An API provides an elegant way to query and mutate the game state. Any client (command line interfaces, browsers) can connect to this API which talks to the backend server, executing the actions needed. There is also a storage that persists every move.
