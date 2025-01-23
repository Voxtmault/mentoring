FROM golang:1.23.5 as build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /pr3

#FROM build-stage AS run-test-stage
#RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /pr3 /pr3

#EXPOSE 47000

USER nonroot:nonroot

ENTRYPOINT ["/pr3"]