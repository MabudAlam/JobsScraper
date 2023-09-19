<img width="1710" alt="Screenshot 2023-09-19 at 3 05 53 PM" src="https://github.com/Pavel401/Jobs-Scraper/assets/47685150/72843e24-8c5c-4faa-bba1-d94e024f53c6">


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

## Configuration

This package relies on a configuration file (`.env`) and a Firebase service account JSON key file (`firebase-config.json`). You should set up these files before running the application.

### `.env` Configuration

Create a `.env` file in the root directory of your project with the following environment variables:

```env
FIREBASE_PROJECT_ID=<your-firebase-project-id>
PORT=<desired-port-number>
UNIVERSAL_DOMAIN=<universal-domain-url>

PRIVATE_KEY_ID=<private-key-id>
PRIVATE_KEY=<private-key>
TYPE=<service-account-type>
CLIENT_EMAIL=<client-email>
CLIENT_ID=<client-id>
AUTH_URI=<auth-uri>
TOKEN_URI=<token-uri>
AUTH_PROVIDER_X509_CERT_URL=<auth-provider-x509-cert-url>
CLIENT_X509_CERT_URL=<client-x509-cert-url>
```

Replace the placeholders with your actual Firebase project details and service account credentials.


## Usage

To run the Jobs Scraper application, simply execute the compiled binary:

```bash
go run main.go
```

The application will start, and you can access the available endpoints in your web browser or via HTTP requests.

## Endpoints

The Jobs Scraper package provides the following endpoints:

- `/cred`: Retrieves the Firebase service account credentials.
- `/atlassian`, `/amazon`, `/coursera`, `/freshworks`, `/gojek`, `/mpl`: Scrapes job postings from specific websites.
- `/all`: Scrapes job postings from all supported websites.
- `/syncFirestore`: Synchronizes job postings with Firebase Firestore.
- `/getJobsFromFirestore`: Retrieves job postings from Firebase Firestore.
- `/`: Serves an HTML file (index.html).

## Contributing

Contributions to this project are welcome! If you have ideas for improvements or encounter any issues, please open an issue or submit a pull request on the [GitHub repository](https://github.com/Pavel401/Jobs-Scraper).
