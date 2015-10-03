# go-api-now

This is a very basic HTTP API that returns current time and environment in Go.

In default, the api server is listening to tcp/8000. You can change the port by specifying `PORT` environment variable.

# Why?

Sometimes I want to test Docker integration so I need a container image suitable for testing.

e.g.

- AWS ECS integration
- AWS CodeDeploy integration
- Kubernetes integration
- Heroku integration
- ...

# Usage

Build:

```sh
$ docker build -t nabeken/go-api-now:latest .
```

Run:

```sh
$ docker run -d -p 8000:8000 nabeken/go-api-now:latest
$ curl -i http://<docker>:8000
```

Or

```sh
$ docker run -d -p 9999:9999 -e PORT=9999 nabeken/go-api-now:latest
$ curl -i http://<docker>:9999
```
