# Database

We are using a PostgreSQL Database for the developement.

### Setting up

You need to setup a PSQL database and provide the connection URL to the same in the `.env` in the root of the package.

### Schema

The schema for the backend API is stored in the [`schema.sql`](./assets/schema.sql) file. You can connect to your remote instance of PostreSQL database and execute the given queries to setup and initalize all tables in the same. You can then populate the said tables with the help of the scrapers.

As an alterantive, if you have got the database locally, then you can use Postgres' backup and migration comand line utilities to move the data to a remote instance, preserving the structures and relations.

You can use the following command from the root to execute the SQL file against the database:

```bash
cat database/assets/schema.sql | psql -h localhost -U admin -d lfx
```

Remember to replace the values of `-h` flag with the host, `-U` with the user and `-d` with the database name. You can even pipe the result directly to `psql <YOUR_POSTGRES_CONNECTION_STRING>`.
