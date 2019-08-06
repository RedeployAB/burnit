FROM alpine:3.10

ARG VERSION

WORKDIR /app

COPY release/bin/secretgen /app/

EXPOSE 3000

ENTRYPOINT [ "./secretgen" ]
