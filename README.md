# URL Shortener

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![SQLite](https://img.shields.io/badge/Database-SQLite-003B57.svg)](https://sqlite.org/)
[![Deployed on Railway](https://img.shields.io/badge/Deployed-Railway-purple.svg)](https://railway.app/)
[![Live Demo](https://img.shields.io/badge/Live-Demo-green.svg)](https://sunny-quietude-production.up.railway.app/)

A minimal production-ready URL shortener built with **Go**, **SQLite**, and deployed on **Railway**.

**Live Demo:**
👉 https://sunny-quietude-production.up.railway.app/

---

## Overview

This project demonstrates:

- Backend development with Go (`net/http`)
- Persistent data storage with SQLite
- Database initialization on application startup
- Anonymous user tracking via cookies
- Production deployment on Railway
- Clean and simple UI

The service converts long URLs into compact, shareable short links.

---

## How It Works

1. User submits a long URL.
2. Server generates a short key (e.g. `360`).
3. Data is stored in SQLite:
   - `short_key`
   - `long_url`
   - `cookie_id`
   - `created_at`
4. When `/wow/{short_key}` is opened:
   - The server queries the database
   - Performs HTTP 302 redirect to the original URL

---

## 🗂 Project Structure

```
URL-Shortener/
│
├── db/
│ └── url_db.db
│
├── static/
│ ├── style.css
│ └── reset.scss
│
├── templates/
│ └── index.html
│
├── main.go
├── go.mod
├── go.sum
└── README.md
```

---

## Database

SQLite database file is created automatically:

Table schema:

```
sql
CREATE TABLE IF NOT EXISTS urls(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_key TEXT NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    cookie_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🍪 Anonymous User Tracking

Each visitor receives a `session_id` cookie.

This allows:

- Identifying links created by the same user
- Tracking without authentication
- Implementing a simple ownership model without user accounts

---

## 🛠 Tech Stack

- **Go 1.24**
- `net/http`
- `database/sql`
- `modernc.org/sqlite`
- HTML Templates
- Railway (Deployment)

---

## Anonymous User Tracking

The application assigns a unique `session_id` cookie to every visitor.

This enables:

- Anonymous link ownership
- Session-based tracking without login
- Lightweight user identification
- Clean and minimal architecture without authentication overhead

## Run Locally

```bash
git clone https://github.com/rrrrr-neko/url-shortener.git
cd url-shortener

go mod tidy
go run .

---

## Author

Built as a backend learning project
by @rrrrr-neko