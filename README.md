# Jobs Scraper

![Screenshot](https://github.com/Pavel401/Jobs-Scraper/assets/47685150/e47e5503-f824-47a2-ad63-ed66c298f350)

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

The Jobs Scraper is a powerful Go application designed for scraping job postings from a variety of websites and storing them in an SQLite database. Additionally, it offers endpoints for retrieving and synchronizing job postings, making it a valuable tool for job seekers and developers alike.

## Prerequisites

Before getting started with the Jobs Scraper, please ensure you have the following prerequisites installed:

- [Go Programming Language](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

## Installation

To install the Jobs Scraper package, follow these simple steps:

1. Clone the repository:
   ```bash
   git clone <repository-url>
   ```

2. Navigate to the project directory:
   ```bash
   cd <repository-directory>
   ```

3. Build the Go application:
   ```bash
   go build
   ```

## Usage

To run the Jobs Scraper application, execute the compiled binary:

```bash
go run main.go
```

This will start the application, allowing you to access its features through your web browser or via HTTP requests.

## Endpoints

The Jobs Scraper package provides a set of powerful endpoints for scraping, managing, and synchronizing job postings. Here's a brief overview of the available endpoints:

- `/atlassian`, `/amazon`, `/coursera`, `/freshworks`, `/gojek`, `/mpl`: Scrapes job postings from specific websites.
- `/syncwithSql`: Synchronizes job postings with a SQL database.
- `/getallJobsFromSQL`: Retrieves job postings from a SQL database.
- `/`: Serves an HTML file (index.html).

## Configuration

To configure the application, you can set environment variables in an `.env` file. For example:

```env
SYNC_WITH_SQL_PASSWORD=password
```

## Contributing

Contributions to this project are highly encouraged! If you have ideas for improvements or come across any issues, please don't hesitate to open an issue or submit a pull request on the [GitHub repository](https://github.com/Pavel401/Jobs-Scraper).

## License

This project is licensed under the [License Name] License - see the [LICENSE](LICENSE) file for details.

---

# To-Do List


1. **Improved Frontend:**
   - [ ] Enhance the user interface of the application.
   - [ ] Implement responsive design for better usability on different devices.
   - [ ] Add styling and improve overall aesthetics.

2. **Filters:**
   - [ ] Implement advanced filtering options for job searches.
   - [ ] Allow users to filter jobs based on location, job type, and other relevant criteria.

3. **Better Scraping Methods:**
   - [ ] Explore and implement more efficient and reliable scraping methods for job postings.
   - [ ] Handle edge cases gracefully and improve the accuracy of data extraction.

4. **Fix Missing Data:**
   - [ ] Address and fix any issues related to missing or incomplete data during the scraping process.
   - [ ] Ensure that all relevant information is captured and stored accurately.

5. **Unit Testing:**
   - [ ] Implement unit tests to ensure the reliability of critical components.
   - [ ] Set up automated testing processes for continuous integration.

6. **Error Handling:**
    - [ ] Improve error handling mechanisms to provide informative error messages to users.
    - [ ] Log errors for debugging and troubleshooting purposes.

7. **Refactoring:**
    - [ ] Refactor code for better readability and maintainability.
    - [ ] Break down complex functions into smaller, more modular components.


