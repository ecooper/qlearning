# qlearning

The qlearning package provides a series of interfaces and utilities to implement
the [Q-Learning](https://en.wikipedia.org/wiki/Q-learning) algorithm in
Go.

This project was largely inspired by [flappybird-qlearning-
bot](https://github.com/chncyhn/flappybird-qlearning-bot).

*Until a release is tagged, qlearning should be considered highly
experimental and mostly a fun toy.*

## Installation

```shell
$ go get https://github.com/ecooper/qlearning
```

## Quickstart

qlearning provides example implementations in the [examples](examples/)
directory of the project.

[hangman.go](examples/hangman.go) provides a naive implementation of
[Hangman](https://en.wikipedia.org/wiki/Hangman_(game)) for use with
qlearning.

```shell
$ cd $GOPATH/src/github.com/ecooper/qlearning/examples
$ go run hangman.go -h
Usage of hangman:
  -debug
        Set debug
  -games int
        Play N games (default 5000000)
  -progress int
        Print progress messages every N games (default 1000)
  -wordlist string
        Path to a wordlist (default "./wordlist.txt")
  -words int
        Use N words from wordlist (default 10000)
```

By default, running [hangman.go](examples/hangman.go) will play millions
of games against a 10,000-word corpus. That's a bit overkill for just
trying out qlearning. You can run it against a smaller number of words
for a few number of games using the `-games` and `-words` flags.

```shell
$ go run hangman.go -words 100 -progress 1000 -games 5000
100 words loaded
1000 games played: 73 WINS 927 LOSSES 7 ACCURACY
2000 games played: 408 WINS 1592 LOSSES 34 ACCURACY
3000 games played: 1108 WINS 1892 LOSSES 70 ACCURACY
4000 games played: 1975 WINS 2025 LOSSES 87 ACCURACY
5000 games played: 2913 WINS 2087 LOSSES 94 ACCURACY

Agent performance: 5000 games played, 2913 WINS 2087 LOSSES 58 ACCURACY
```

Accuracy per progress report is isolated within that cycle, a 5000 games
in this example. The accuracy report is meant to show the velocity of
learning by the agent. The accuracy itself is mostly irrelevant at this
point.

As you can see, after 5000 games, the agent is able to "learn" and play
hangman against a 100-word vocabulary.

## Usage

See [godocs](https://godoc.org/github.com/ecooper/qlearning) for the
package documentation.
