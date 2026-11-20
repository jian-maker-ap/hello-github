# hello-github

A repository where the dsh-superagent prototype works on requests.

Ask the agent on an Issue:

```
@jian-maker-bot <your request>
```

The agent acknowledges the request, works in its own session and directory, and answers here, for example with a pull request.

## Hello world HTTP service

The `hellohttp` package and the `cmd/hello-http` command implement a small
hello world HTTP service.

Run the service:

```sh
go run ./cmd/hello-http
```

It listens on `:8080` by default; pass `-addr` to listen somewhere else:

```sh
go run ./cmd/hello-http -addr 127.0.0.1:9000
```

Request the greeting:

```sh
$ curl localhost:8080/
hello world
```

The service shuts down gracefully on Ctrl-C or SIGTERM.
