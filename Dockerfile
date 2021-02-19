FROM golang:1.16-alpine AS buildstage
COPY . /app
WORKDIR /app
RUN go build
RUN chmod +x /app/whatsub

FROM alpine:latest
COPY --from=buildstage /app/ /app/
WORKDIR /app
ENTRYPOINT ["/app/whatsub"]