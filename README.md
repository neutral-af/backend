# carbon-offset-backend

# Getting Started

First, install the Go dependencies and build the GraphQL generated types (if necessary):

    make deps

To run unit tests, use:

    make test

To run the dev server, export a `CLOVERLY_API_KEY`, then run:

    make dev

This dev server starts an http dev server on port 8080, equivalent to the http routing set up by Zeit.

To simulate the runtime of Zeit, after installing the `now` cli, run: (you'll need a `.env` file with the necessary environment variables)

    now dev

## GraphQL Schema

The schema is imported via a git submodule at `schema/`, and the Go types are generated using `gqlgen` and configured in the `gqlgen.yml`.

When the schema changes, the generated Go types need to be updated as well. After updating the `schema` submodule, run:

    make deps


