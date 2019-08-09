FROM alpine:3.10

ARG VERSION
ARG PORT=4000

WORKDIR /app

COPY release/bin/secretdb /app/

EXPOSE ${PORT}

ENTRYPOINT [ "./secretdb" ]
