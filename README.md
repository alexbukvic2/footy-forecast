# footy-forecast

A prediction game for major football tournaments (starting with the World Cup).
Users predict match results, top scorers, group winners, and the tournament champion.
Players can join private leagues and compete on global and per-league leaderboards.

## Stack
- **Backend:** Go 1.26
- **Database:** PostgreSQL (RDS)
- **Auth:** AWS Cognito
- **Hosting:** AWS EC2 (single instance, free tier)
- **CI/CD:** GitHub Actions

## Local development

Requirements:
- Go 1.26+
- golangci-lint

```bash
make run     # start the server on :8080
make test    # run tests
make lint    # run linters
```

Health check:
```bash
curl http://localhost:8080/health
```

## Project layout

See [docs/decisions/0001-project-layout.md](docs/decisions/0001-project-layout.md).
