# Introduction

This is a sample integration test suite for a database access interface written in golang.

Traditionally, database access testing in unit tests is accomplished using mocks. This test suite demos a full integration test that runs by spinning up a docker instance of postgres, connecting to that instance, running a set of table-driven tests against that database, and deleting the container once the tests have completed.

This suite can be extended to include test cases for multiple databases (i.e. MySQL, Mongo etc) because the concrete implementation of the database access is abstracted using an interface.

Finally, since an actual database instance is being brought up we can test features like database migrations against a real database along with our unit tests, which would otherwise not be possible via mocks.

## Key Highlights

- Database access is declared in an interface found in `db.go` allowing for multiple concrete implementations (however only the postgres one has been implemented in this example)
- Database access for the suite & connection instances are provided to the suite using dependency injection
- Test suite is structured in a format that is more common to frameworks like Mocha & Jasmine (setup, teardown, etc.), which is not strictly speaking "idiomatic" go code, but it makes for a much easier time in terms of maintenance & enhancement
- Table driven tests are run against a real database, bringing local dev environment closer to prod

## Required Dependencies

- Docker (recommended to download the `postgres:latest` image otherwise the first time running the suite will take a lot longer)
- Golang runtime in path

## Run Instructions

1. Start docker daeomon (must be running for tests to run)
2. Run `go test ./... -v`
