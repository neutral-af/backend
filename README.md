# carbon-offset-backend

# Getting Started

First, install the Go dependencies and build the GraphQL generated types (if necessary):

    make deps

To run unit tests, use:

    make test

To run the dev server, export a `CLOVERLY_API_KEY`, then run:

    make dev

This dev server starts an http dev server on port 8000, equivalent to the http routing set up by Zeit.

To simulate the runtime of Zeit, after installing the `now` cli, run: (you'll need a `.env` file with the necessary environment variables)

    now dev

## GraphQL Schema

The schema files are located in `schema/`, and the Go types are generated using `gqlgen` and configured in the `gqlgen.yml`.

When the schema changes, the generated Go types need to be updated as well. After updating the `schema` submodule, run:

    make deps

## Config & Credentials

The environment variables required are specified in [config.go](lib/config/config.go). They need to be either exported in the environment or provided in a dotenv file.

## Example Queries

Health check:

    { health }

Estimate for a single flight:

    {
        estimate {
            fromFlights(flights:{ departure:"YYZ", arrival:"LHR" }) {
                carbon price { currency, cents }
            }
        }
    }

Estimate from multiple flights, overriding price currency:

    {
        estimate {
            fromFlights(flights:[
                { departure:"YYZ",arrival:"LHR" },
                { departure:"LHR",arrival:"YYZ" }
            ]) {
                price(currency:EUR) { currency, cents } carbon provider
            }
        }
    }
