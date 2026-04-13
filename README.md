<div align="center">
  <img src="./Docs/image/banner.png" alt="Door - API Gateway" width="700"/>

  <h1>Door</h1>
  <p>A lightweight, high-performance API Gateway written in Go.</p>

  ![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
  ![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)
  ![Angular](https://img.shields.io/badge/Angular-21-DD0031?style=flat&logo=angular)
  ![License](https://img.shields.io/badge/license-MIT-green?style=flat)
</div>

---

**Door** sits between your clients and your backend services — handling authentication, routing, and logging so your APIs don't have to.

## Features

- **Brute force protection** — protection for brute force
- **Request Logging** — every request is captured asynchronously: method, path, IP, response time, status code, and more
- **CORS Middleware** — configurable per-origin with preflight support
- **Dashboard** — built-in frontend for managing applications and auth flows


## Tech Stack

| Layer | Technology |
|---|---|
| Core | Golang |
| Data | SQLite3 & Redis |
| View | Angular |

## API Reference

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive tokens |
| `GET` | `/api/v1/auth/token` | Validate / refresh access token |
| `POST` | `/api/v1/auth/logout` | Logout and clear session |

### Applications (proxy targets)

All routes below require `Authorization: Bearer <access_token>`.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/applications` | List all registered applications |
| `GET` | `/api/v1/applications/:id` | Get a single application |
| `POST` | `/api/v1/applications` | Register a new application |
| `DELETE` | `/api/v1/applications/:id` | Remove an application |

### Example: Add an application

```bash
curl -X POST http://localhost:7171/api/v1/applications \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"url": "http://my-service:8080", "country": "Brazil"}'
```

## Token Strategy

- **Access token** — short-lived (15 min), kept in memory on the client
- **Refresh token** — long-lived (24h), stored in localStorage
- Protected routes validate the `Authorization: Bearer` header via middleware

## License

MIT
