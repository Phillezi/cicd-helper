FROM --platform=$BUILDPLATFORM golang:bookworm AS builder

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update &&\
    apt-get install -y make

COPY . .

RUN make

FROM alpine:latest AS runner

COPY --from=builder /app/bin/cicd-helper /bin/cicd-helper

EXPOSE 8080

CMD [ "./bin/cicd-helper" ]
