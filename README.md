# Testing OpenID Connect and OAuth 2.0 with WireMock & Testcontainers - use case from Dynatrace
## Overview
This repository contains demo code that was prepared for a presentation on the use case of open source projects such as WireMock and Testcontainers for testing puroses within Dynatrace. The code demonstrates how to use OpenID Connect and OAuth in Go for authentication and authorization purposes. Specifically, it shows the process of calling the OpenID Connect discovery endpoint to fetch server configurations and then using the token_endpoint value from the response to make a second call to the /token endpoint. This retrieves an access token using the client credentials flow. This setup is designed to run tests in parallel by default, utilizing the `t.Parallel()` call at the beginning of tests.
###
![overview_paralel](https://github.com/AntoniKania/oauth-wiremock-testcontainers-go/assets/87483058/f51a5e60-1829-49c4-962c-db0b65f9556c)

## Key Features
**OpenID Connect Discovery:** Automated fetching of server configuration details.
**OAuth Client Credentials Flow:** Secure method to obtain access tokens for microservices.
**Parallel Test Execution:** Speed up testing with Go's parallel execution feature.
**WireMock Integration:** Simulates the authorization server for testing, providing both the well-known and token endpoints.
**Testcontainers Support:** Manages the WireMock server instance within a Docker container directly from the test code.


## Running locally
### Prerequisites
**Go:** You must have Go installed on your machine. Visit Go's official website for installation instructions.
**Docker:** Required for running the WireMock server instance inside a container. Download Docker from Docker's official website.
### Run the tests:

This project is configured to run tests in parallel by default, significantly reducing the test execution time. To run the tests, use the command:

```bash
go test ./...
```
This will automatically start the WireMock server in a Docker container, execute all tests, and then shut down the server.

![sequence_diagram](https://github.com/AntoniKania/oauth-wiremock-testcontainers-go/assets/87483058/ddc7da29-f611-4c99-94d1-95c8754be619)


## Presentation
For a detailed explanation of the code and its use case within Dynatrace, check out the presentation! :)
