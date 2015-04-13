# httpq

With httpq, you can buffer HTTP requests and replay them later, and either Redis or BoltDB can be used for persistence. This is useful for buffering HTTP requests that do not have to be processed in realtime, such as webhooks.

## Installation

```shell
$ go install github.com/DavidHuie/httpq/cmd/httpq
```

## Using Redis

```shell
$ httpq -redis=true -redis_url=":6379"
```

## Using BoltDB (disk persistence)

```shell
$ httpq -db_path="/tmp/httpq.db"
```

## Queuing a request

```shell
$ curl localhost:3000/push
```

## Replaying a request

```shell
$ curl localhost:3000/pop

GET /push HTTP/1.1
Host: localhost:3000
Accept: */*
User-Agent: curl/7.37.1
```

## Determining size of queue

The result is returning as JSON.

```shell
$ curl localhost:3000/size

{"size":3}
```
