FROM golang:1.25.1@sha256:a5e935dbd8bc3a5ea24388e376388c9a69b40628b6788a81658a801abbec8f2e AS build

WORKDIR /build

COPY . .

RUN go build -o /proxy ./cmd/proxy

FROM gcr.io/distroless/base-debian12@sha256:d605e138bb398428779e5ab490a6bbeeabfd2551bd919578b1044718e5c30798

COPY --from=build /proxy /proxy

ENTRYPOINT ["/proxy"]
