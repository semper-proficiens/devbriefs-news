# devbriefs-news
Devbriefs microservice about news

# Go

## News-Fetcher service

Workflow:
- `config`: app configuration
- `models`: how data will be structure to handle data from API
- `handlers`: all api handlers for our service
- `service`: business logic
- `utils`: any tools we develop to help our components

## Design Principles

### Repository Pattern

To separate the part of the application that deals with data (like fetching from API or DB) is separated from the
rest of the application (business logic). We can test data logic and business logic separately and easier. But, the
idea is that the data source (News API can be switched effortlessly).

- e.g. For example, fetching news articles from an API

In our implementation that's the `internal/api` components.

### Service Layer Pattern

To separate business logic from user interface and data code. Meant to reuse some parts of this logic in other `devbriefs`
service components.

- e.g. For example, deciding which news articles to fetch based on user input and ensuring all rules are followed

In our implementation that's the `internal/service` components.

# Testing

## Local simple test

We can run our service locally like this:
```go
NEWSFETCHER_GOOGLE_API_KEY=$apiKey go run cmd/news-fetcher/main.go
```

We can query the NewsAPI directly in simple curl like this:

Using `everything` endpoint:
```bash
curl https://newsapi.org/v2/everything -G \
    -d q=Apple \
    -d from=2024-08-25 \  
    -d sortBy=popularity \
    -d apiKey=$apiKey
```

Using `top-headlines` endpoint:
```bash
curl https://newsapi.org/v2/top-headlines -G \
    -d q=Apple \
    -d country=us -d category=technology \
    -d apiKey=$apiKey
```

Using our local api endpoint for everything news:
```bash
 curl -X GET "http://localhost:8080/api/everything-hacking-news"
```

## Go Tests

```bash
 make go_tests
```

# TOIL

- Grab `companies` or `domains` of interest from the request query to this api, and construct another hacking news function
targeting those only
- Add unit, integration and load test
- Improve app performance with go routines, and fan out
- Add automatic linters in CI
- Setup for quality scores, code smells, etc
- Explore other protocols like gRPC
- Store news in backend and add caching for it


