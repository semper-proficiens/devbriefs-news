# devbriefs-news
Devbriefs microservice about news

# Go

## Design Principles

### Repository Pattern

To separate the part of the application that deals with data (like fetching from API or DB) is separated from the
rest of the application (business logic). We can test data logic and business logic separately and easier. But, the
idea is that the data source (News API can be switched effortlessly).

- e.g. For example, fetching news articles from an API

### Service Layer Pattern

To separate business logic from user interface and data code. Meant to reuse some parts of this logic in other `devbriefs`
service components.

- e.g. For example, deciding which news articles to fetch based on user input and ensuring all rules are followed