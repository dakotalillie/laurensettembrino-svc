.Phony: build-SendEmailFunction

build-SendEmailFunction:
	GOARCH=amd64 GOOS=linux go build -o bootstrap . \
	&& mv bootstrap $(ARTIFACTS_DIR)
