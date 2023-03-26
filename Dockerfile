# syntax=docker/dockerfile:experimental
# ---
FROM golang:1.20 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

# Build admission-webhook
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
  go build -o bin/admission-webhook .


# ---
FROM alpine AS run

COPY --from=build /work/bin/admission-webhook /usr/local/bin/
#COPY cli_tools/syft /usr/local/bin
# add the syft and grype dependencies
RUN apk --no-cache add curl
RUN curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
RUN curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b /usr/local/bin
ENV PATH="${PATH}:/usr/local/bin"

COPY pkg /usr/webhook
WORKDIR /usr/webhook

CMD ["admission-webhook","-hook"]
