# API Endpoints

## Job Scraper API

Base URL: `http://localhost:8080`

---

## GET /

Health check endpoint.

**Response:**
```json
{
  "status": "Job Scraper API"
}
```

---

## GET /syncall

Runs the full scraping pipeline for all enabled job providers (e.g., AshbyHQ).

**Query Parameters:**
| Parameter | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| password  | string | Yes      | Password from `SYNC_WITH_SQL_PASSWORD` env var |

**Example:**
```
GET /syncall?password=your_password_here
```

**Response (200 OK):**
```json
{
  "status": "sync completed"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid password"
}
```

---

## GET /getallJobsFromSQL

Returns all active jobs stored in the SQLite database.

**Response (200 OK):**
```json
{
  "jobs": [
    {
      "jobName": "Senior Software Engineer",
      "description": "We are looking for...",
      "date": "2026-04-10T10:00:00Z",
      "applyLink": "https://apply.example.com/job/123",
      "companyName": "Example Co",
      "meta": {
        "location": "San Francisco",
        "remote": true,
        "department": "Engineering",
        "team": "Platform",
        "employmentType": "Full-time",
        "compensation": "$150k - $200k",
        "jobUrl": "https://example.com/jobs/123",
        "contentHash": "abc123def456",
        "source": "ashby"
      }
    }
  ]
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "error message here"
}
```

---

## Environment Variables

| Variable              | Required | Default    | Description                          |
|-----------------------|----------|------------|--------------------------------------|
| `SYNC_WITH_SQL_PASSWORD` | Yes      | -          | Password for `/syncall` endpoint      |
| `DATABASE_PATH`       | No       | `./jobs.db`| Path to SQLite database file         |
| `RAILWAY_ENVIRONMENT` | No       | -          | Set to skip `.env` loading           |
| `DEBUG`               | No       | -          | Set to `1` to enable debug logging   |
| `FETCH_DELAY_MIN`     | No       | 100        | Min delay between company scrapes (ms)|
| `FETCH_DELAY_MAX`     | No       | 500        | Max delay between company scrapes (ms)|
| `MIN_SCORE`           | No       | 10         | Min relevance score for notifications|

### Company Configuration

Configure which companies to scrape:

**Option 1: JSON array (ASHBY_COMPANIES)**
```bash
ASHBY_COMPANIES='[{"Company":"Vercel","AshbySlug":"vercel","Enabled":true},{"Company":"Linear","AshbySlug":"linear","Enabled":true}]'
```

**Option 2: Comma-separated (ASHBY_COMPANIES_COMMA)**
```bash
ASHBY_COMPANIES_COMMA="Vercel:vercel,Linear:linear,Notion:notion"
```

**Default if none set:** Vercel

---

## Future Endpoints

When adding new scrapers, routes will be prefixed by provider:

| Endpoint                  | Description                              |
|---------------------------|------------------------------------------|
| `GET /syncall?password=X` | Syncs all providers (current)            |
| `GET /getallJobsFromSQL`  | Gets all jobs from database (current)     |
| `GET /sync/ashby`         | Sync only AshbyHQ (future)               |
| `GET /syncgreenhouse`     | Sync Greenhouse jobs (future)            |
| `GET /jobs`               | Query/filter jobs (future)               |
