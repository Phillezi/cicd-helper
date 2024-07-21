FROM --platform=$BUILDPLATFORM golang:bookworm AS builder

WORKDIR /app

COPY . .

RUN make

FROM alpine:latest AS runner

COPY --from=builder /app/bin/cicd-helper .

EXPOSE 8080

ENTRYPOINT [ "./cicd-helper" ]
