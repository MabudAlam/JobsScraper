<img width="1470" alt="Screenshot 2023-09-28 at 2 58 49 AM" src="https://github.com/Pavel401/Jobs-Scraper/assets/47685150/e47e5503-f824-47a2-ad63-ed66c298f350">


# Jobs Scraper 

This is the README file for the Jobs Scraper package, a Go application that scrapes job postings from various websites and stores them in a Firebase Firestore database. This package also provides endpoints to retrieve and synchronize job postings.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Contributing](#contributing)
- [License](#license)

## Introduction

The Jobs Scraper package is designed to scrape job postings from specific websites, including Atlassian, Amazon, Coursera, FreshWorks, Gojek, and MPL. It stores the scraped job data in a Firebase Firestore database, making it easy to access and manage job postings.

## Prerequisites

Before you can use this package, ensure that you have the following prerequisites installed:

- [Go Programming Language](https://golang.org/)
- [Firebase Firestore](https://firebase.google.com/docs/firestore) (with a service account JSON key)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Viper](https://github.com/spf13/viper)
- [goccy/go-json](https://github.com/goccy/go-json) (for JSON serialization)

## Installation

To install the Jobs Scraper package, you need to clone the repository and build the Go application.

```bash
git clone <repository-url>
cd <repository-directory>
go build
```

## Usage

To run the Jobs Scraper application, simply execute the compiled binary:

```bash
go run main.go
```

The application will start, and you can access the available endpoints in your web browser or via HTTP requests.

## Endpoints

The Jobs Scraper package provides the following endpoints:
Certainly! Here's the updated README for your Jobs Scraper package with the provided code changes:

---

# Jobs Scraper Package

The Jobs Scraper package allows you to scrape job postings from various websites and manage the data in Firebase Firestore. This package provides the following endpoints:

- `/cred`: Retrieves the Firebase service account credentials.
- `/atlassian`, `/amazon`, `/coursera`, `/freshworks`, `/gojek`, `/mpl`: Scrapes job postings from specific websites.
- `/syncwithSql`: Synchronizes job postings with a SQL database.
- `/getallJobsFromSQL`: Retrieves job postings from a SQL database.
- `/`: Serves an HTML file (index.html).



## Contributing

Contributions to this project are welcome! If you have ideas for improvements or encounter any issues, please open an issue or submit a pull request on the [GitHub repository](https://github.com/Pavel401/Jobs-Scraper).
