# devbriefs-news
Devbriefs is microservice about catered Engineering news

[//]: # (<img align="right" width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">)

[![Build Status](https://github.com/semper-proficiens/devbriefs-news/actions/workflows/gotests.yml/badge.svg)](https://github.com/semper-proficiens/devbriefs-news/actions?query=branch%3Amain+)
[![codecov](https://codecov.io/github/semper-proficiens/devbriefs-news/branch/main/graph/badge.svg?token=75SCUZRRIP)](https://codecov.io/github/semper-proficiens/devbriefs-news)

## DevBriefs News service

Overall Concept:
- We fetch news from a News API
- To avoid making subsequent API calls and get the objects faster we cache the news
- Every time someone hits an API endpoint, it fetches news and updates cache
- There is a daily routine that run every day 6am EST, after executing it will sleep again until next day same time
- All cached news will have a ttl of 24 hours
- There is a unique check to only insert unique titles based on word similarity in the titles
- The titles are hashed for uniqueness based on article title

Repo Structure:
- `api`: 3rd party apis
- `datastore`: our backends and caches
- `handlers`: all api handlers for our service
- `models`: json models expected from certain 3rd party apis
- `service`: business logic

# Testing

## Local simple test

We can run our service locally like this (**needs to be a valid api key**):
```go
GOOGLE_NEWS_API_KEY=$apiKey go run main.go
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

## Go Tests and Lints

To run local go tests, with benchmarks, coverage, lints, vets, and gosec:
```bash
 make go_tests
```

1. Install golangci-lint https://golangci-lint.run/welcome/install/#local-installation
2. Run make command:
    ```bash
    make golangci_run
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


