# LFX Mentorship Metrics

`LFX Mentorship Metrics` is your one-stop solution for gaining valuable insights into the world of open source project mentorship within the Linux Foundation's LFX ecosystem. This backend API serves is a microservice built in Go, which parses that data related to the Mentorship, generates statistics from them, and exposes the same via a REST API.

## Need

The `LFX Mentorship` data is highly unstructed. There is no way to use the present [`mentorship.lfx site`](https://mentorship.lfx.linuxfoundation.org/) to view the data of organizations and projects across years, and evaluate trends so that new contributors can figure out which organizations and projects would be the most relevant to them.

## Architecture

The `scrapers` package is meant to contain all the scripts that can be run on a periodic basis to collect all the data related to the LFX mentorship projects in a SQL database. The `api` package contains all the source code for the controllers and endpoints exposed, which are used by the frontend to render the `LFX Mentorship Metrics` project.

## Setting Up the Development Environment

As the API relies on the connection to a PostgreSQL database, we will be using Docker containers to manage the same and populate it with data for a uniform development experience. Follow the following steps for the same:

1. Install [Docker](https://docs.docker.com/engine/install/) and [Docker Compose](https://docs.docker.com/compose/install/).
2. In command line, run `docker compose up`.

The API is live now.

### API

The API would now be served on `0.0.0.0:8080`. You can go to the `0.0.0.0:8000/api` to get a health check for the same.

As `Go` does not supply any HMR package, after making changes to the source code, you would need to close the container from the terminal (use `Ctrl + C`) and then restart the same by using `docker compose up`. This would be a quick process, as Docker caches all the layers.

### Database

The PostgreSQL server would be exposed at `localhost:8079`.

From this project, we have configured a user `admin` (with the super secure password `admin`) having all the previligies on the database `lfx` (which contains all our data). To interact with the database, you can connect to the same from your terminal using:

```bash
psql -U admin -h localhost -p 8079 -d lfx
```

and entering the password as `admin` when prompted. This requires you to have the `psql` CLI installed on your system.

If you make any changes to the database, then you can use the `pg_dump` utility to store the same in the [`backup.sql`](./sql/backup.sql), which is used to populate the database intially. Doing so would allow all other developers to work on the updated database. To do so, you can just run the following command after making all the necesary changes to the database:

```bash
pg_dump -U admin -h localhost -p 8079 lfx  > ./sql/backup.sql
```

#### Debug Help

If you have changed the `backup.sql` file and the new data is not being reflected in your database, then it you might need to delete all the Docker's cache and build volumes by running the command `docker system prune -a`, and then re-running the container with `docker compose up`.

### Individual Docker Images

You can play around with by builing the images for the API and database separately if you want more fine grained control. Be sure to provide the relevant inputs and environment configurations for the images as done by the [`docker compose file`](./docker-compose.yml).

```bash
docker build -t db . -f sql.Dockerfile
docker build -t api . -f api.Dockerfile
```
