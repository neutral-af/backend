# Neutral.af - API

*The Golang and GraphQL project that's carbon neutral (as f\*\*k).*

## Contributing

All contributions are welcome! Please feel free to open issues or PRs, update documentation, or implement features. Look for issues marked as `good-first-issue` for good places to jump in.

If you're considering implementing a new feature or modifying an existing one, please open an issue first for discussion (so that we can talk about implementation before you spend your precious time on it).

## Developing

### Quick Start

To run the dev server, copy the `.env.template` to a `.env` file. Add the necessary keys (see [below](#config--credentials)) and then run:

    make dev

This dev server starts an http dev server on port 8000, matching the http routing set up in production/staging by [Zeit](https://zeit.co).

A [GraphQL playground](https://github.com/prisma-labs/graphql-playground) is also found on `/playground` (ie. http://localhost:8000/playground). You can try entering the query `{ health { aliveSince environment } }` to test the playground's connection to the backend.

To run the Go unit tests, use:

    make test

### Prerequisites

The only core requirement is a recent version of Go, with support for Go Modules.

If you want to deploy your own test version to Zeit, you'll need the [now-cli](https://github.com/zeit/now).

If you want to change the GraphQL schema, you'll also need to install [gqlgen](https://github.com/99designs/gqlgen), with `go get -u github.com/99designs/gqlgen`.

### GraphQL Schema

The schema files are located in `schema/`, and the Go types are generated using `gqlgen` and configured in the `gqlgen.yml`. You'll need to [install gqlgen](https://github.com/99designs/gqlgen) first.

When the schema changes, the generated Go types need to be updated as well. After updating any of the `.graphql` files, run:

    make deps

### Config & Credentials

The environment variables required are specified in [config.go](lib/config/config.go). They need to be either exported in the environment or provided in a dotenv (`.env`) file. If using environment variables, `DEV_` must be added as a prefix. A `.env.template` file is provided as an example.

The keys you need will depend on what you're working on:

- If you are only developing the carbon-estimation logic, you won't need any keys (but you won't be able to make any GraphQL queries, so you'll need to use on tests)
- If you're working on the logic for calling one of the offset providers (eg. Cloverly), you'll need a test/sandbox API key (which you can typically get by signing up for a test account with the provider)
- If you are working on the payments logic, you'll need a Stripe private key, which you can get by signing up for a test/sandbox account (if you're also testing the frontend, you'll need the matching public key)

## Example Queries

Health check:

    { health { aliveSince environment } }

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
                price(currency:EUR) { currency, cents } id carbon provider
            }
        }
    }

Estimate showing fees (via cost breakdown):

    {
        estimate {
            fromFlights(flights:[
                {departure:"YYZ",arrival:"LHR"}
            ]) {
                carbon
                provider
                price(currency:EUR) { currency cents breakdown {name cents currency} }
            }
        }
    }

Purchase an offset using an estimate ID and provider:

    mutation {
        purchase {
            fromEstimate(estimateID:"20190710-cloverly-slug", provider:Cloverly) {
                id carbon details
            }
        }
    }

## Manual Deploy

Zeit doesn't yet provide a sufficient "staging environment", that deploys for master but uses non-production credentials. The frontend system relies on a stable non-production target for its preview environments, so this project can be deployed in a manual non-production mode. To release this, simply run:

    now --env NOW_GITHUB_COMMIT_REF=local

This will make a non-production "staging" deploy, and print its URL (eg. https://backend-jasongwartz.neutral-af.now.sh). This URL can then be configured for the frontend application.
