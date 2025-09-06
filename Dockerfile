FROM golang:1.25 AS build

COPY . .

RUN go build -o /proxy ./cmd/proxy

FROM gcr.io/distroless/base-debian12

COPY --from=build /proxy /proxy

ENTRYPOINT ["/proxy"]
