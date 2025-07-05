# go-api-now

`go-api-now` is a very basic HTTP API that returns the current time written in Go.

In default, the api server is listening to tcp/8000. You can change the port by specifying `PORT` environment variable.

# Why?

Sometimes I want a simple deployable HTTP server to test various integrations.

For example:
- To test Amazon ECS integration
- To test Kubernetes integration
- To test a release pipeline
- To test CI/CD pipeline
- ...

# Endpoints

- `GET /api/ping`
  - returns the current environment variables
  - `?sleep=<duration>`: sleep [`<duration>`](https://golang.org/pkg/time/#ParseDuration) before returns a response
- `GET /api/json`
  - returns the large JSON file
- `GET /api/_stats`
  - returns the current Go's runtime stats
- `GET /api/events`
  - streams the current time with Server-Sent Events ("SSE")

# Usage

## Frontend (Next.js)

Development server:
```sh
npm run dev
```

Production build:
```sh
npm run build
npm run start
```

## Backend (Docker)

Build:
```sh
docker compose build
```

Run:
```sh
docker compose up -d
curl -i http://127.0.0.1:8000
```
