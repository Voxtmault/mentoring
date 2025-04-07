# Using 1 Dockerfile for all services
# Since each service might use the 
FROM golang:1.23.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy shared code / dependencies
COPY pkg/ ./pkg/

ARG SERVICE

# Copy the related service code
COPY internal/${SERVICE}/*.go ./service/
COPY cmd/${SERVICE}/*.go ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -o /${SERVICE}-service

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

ARG SERVICE
COPY --from=build-stage /${SERVICE}-service /${SERVICE}-service

USER nonroot:nonroot

ENTRYPOINT ["/${SERVICE}-service"]