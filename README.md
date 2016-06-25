# ori
[![GoDoc](https://godoc.org/github.com/the-information/ori?status.svg)](https://godoc.org/github.com/the-information/ori)
[![CircleCI](https://circleci.com/gh/the-information/ori.svg?style=svg)](https://circleci.com/gh/the-information/ori)

Write massively scalable, low-latency REST/JSON APIs in Go for App Engine, and have fun doing it.

## what's this thing?

[Google App Engine](https://cloud.google.com/appengine/docs/go/) is a treasure trove of scalable,
high-performance computing that costs very little money.
It just has a couple of minor inconveniences that tend to befuddle people interested in setting up a 5-
minute API:

- No environment variables for app configuration, such as third-party API secrets and the like.
- A very weird and hard-to-use account system.
- No built-in support for REST/JSON APIs unless you want to use Google Cloud Endpoints, which, let's face it,
most of us don't.

ori works with [kami](https://github.com/guregu/kami) to solve those problems,
so you can use App Engine's indefinite stockpile of high-performance computing and storage
in your Go application with very little initial investment of effort.

Read the documentation on [GoDoc](https://godoc.org/github.com/the-information/ori).

## getting started

If you need a guided tutorial, [we have one.](https://github.com/the-information/ori/blob/master/tutorial/01-getting-started.md).

## why's it called ori?

It's a terrible joke having to do with folding paper deities. 
