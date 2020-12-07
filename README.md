# laurensettembrino-svc

This repository contains the source code and infrastructure configuration for backend services used
by laurensettembrino.com. At the moment, this is just a single service for sending
emails submitted via the contact form.

## Requirements

- AWS CLI already configured with Administrator permission
- [Docker installed](https://www.docker.com/community-edition)
- [Golang](https://golang.org)
- SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Local Development

### Build & Compile

```sh
make
```

### Run Tests

Tests need to be run from from the Go module root

```sh
cd send-email && go test
```

### Run a Local API

You can invoke the function at `http://localhost:3000/send-email`. Function contents can be refreshed by running `make`.

```bash
sam local start-api
```

### Deploy

```bash
sam deploy
```
