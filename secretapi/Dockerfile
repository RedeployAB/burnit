FROM alpine:3.10

ARG VERSION
ARG PORT=3000

WORKDIR /app

COPY release/bin/secretapi /app/

EXPOSE ${PORT}

ENTRYPOINT [ "./secretapi" ]
