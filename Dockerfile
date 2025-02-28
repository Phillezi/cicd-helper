FROM golang:bookworm AS builder

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y make

COPY . .

RUN make

FROM debian:stable-slim AS runner

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/cicd-helper /bin/cicd-helper

EXPOSE 8080

CMD [ "/bin/cicd-helper" ]
