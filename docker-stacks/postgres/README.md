# postgres

Local Postgres 17 stack for development, with an optional `pgadmin` profile for a browser UI.

## What's inside

| Service | Image | Default port | Default credentials |
|---|---|---|---|
| `postgres` | `postgres:17-alpine` | `5432` | `piltover` / `piltover` / `piltover` (user / password / db) |
| `pgadmin` (profile `ui`) | `dpage/pgadmin4:latest` | `5050` | `admin@piltover.local` / `admin` |

Healthcheck on the postgres service uses `pg_isready`. `pgadmin` waits for it to be healthy before starting.

## Start it

```bash
piltover stacks up postgres
```

To also start the pgAdmin UI:

```bash
docker compose -f docker-stacks/postgres/compose.yaml --profile ui up -d
```

## Connect

From your host:

```
postgres://piltover:piltover@localhost:5432/piltover
```

From another container on the same Docker network, use `postgres` as the host instead of `localhost`.

## Reset

```bash
piltover stacks nuke postgres   # docker compose down -v — wipes pgdata volume
```

## Customise

Copy `.env.example` to `.env` and edit. Variables: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `PGADMIN_EMAIL`, `PGADMIN_PASSWORD`, `PGADMIN_PORT`.
