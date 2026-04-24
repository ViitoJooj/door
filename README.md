<div align="center">
  <img src="./Docs/image/banner.png" alt="ward - API Gateway" width="900"/>

  <h1>ward</h1>
  <p>Modular API Gateway in Go focused on security, reliability, and observability.</p>

  ![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
  ![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)
  ![Angular](https://img.shields.io/badge/Angular-21-DD0031?style=flat&logo=angular)
  ![License](https://img.shields.io/badge/license-MIT-green?style=flat)
</div>

---

## What Ward is

Ward is a modular API Gateway written in Go.  
It sits between clients and backend services to centralize authentication, request controls, and request logging.

In multi-service architectures, these concerns often become inconsistent across services. Ward provides a single boundary layer so internal services can stay focused on business logic.

## How it works under the hood

A request follows a predictable pipeline:
1. **CORS middleware** validates origin and preflight behavior.
2. **Authentication middleware** validates `access_token` on protected routes.
3. **Handlers** parse/format HTTP only.
4. **Services** execute business rules.
5. **Repositories** perform persistence operations.
6. **Asynchronous logging** writes request metadata without blocking the response.

This separation keeps boundary concerns, domain logic, and data access isolated and maintainable.

## Architecture

- **Handlers**: HTTP interface.
- **Services**: business orchestration and rules.
- **Repositories**: data access layer.
- **DTOs**: input/output API contracts.

## Security posture

- Passwords are hashed with bcrypt before storage.
- Access and refresh tokens are signed and validated by token type.
- Protected endpoints validate identity before business logic runs.
- Access tokens use short expiration windows.
- Request metadata is stored for auditing and incident investigation.

## Reliability and performance

- Built in Go using `fasthttp` for low overhead.
- Request logging is asynchronous to preserve latency.
- Layer boundaries reduce side effects between modules.
- Security controls are centralized at the gateway layer.

## Documentation

- Portuguese version: [`./Docs/README-PT-BR.md`](./Docs/README-PT-BR.md)
- OpenAPI/Swagger: [`./Docs/swagger/core-api.yml`](./Docs/swagger/core-api.yml)
- Contributing guide: [`./CONTRIBUTING.md`](./CONTRIBUTING.md)
- Security policy: [`./SECURITY.md`](./SECURITY.md)
- License: [`./LICENSE`](./LICENSE)

## License

MIT
