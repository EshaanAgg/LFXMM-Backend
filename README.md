# LFX Mentorship Metrics

`LFX Mentorship Metrics` is your one-stop solution for gaining valuable insights into the world of open source project mentorship within the Linux Foundation's LFX ecosystem. This backend API serves is a microservice built in Go, which parses that data related to the Mentorship, generates statistics from them, and exposes the same via a REST API.

## Need

The `LFX Mentorship` data is hughly unstructed. There is no way to use the present [`mentorship.lfx site`](https://mentorship.lfx.linuxfoundation.org/) to vie the data of organizations and projects across years, and evaluate trends so that new contributors can figure out which organizations and projects would be the most relevant to them.

## Architecture

The `scrapers` package is meant to contain all the scripts that can be run on a periodic basis to collect all the data related to the LFX mentorship projects in a SQL database. The `api` package contains all the source code for the controllers and endpoints exposed, which are used by the frontend to render the `LFX Mentorship Metrics` project.
