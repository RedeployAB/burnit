FROM scratch

ARG VERSION
ARG PORT=3000

WORKDIR /app

COPY release/bin/secretgw /app/

EXPOSE ${PORT}

ENTRYPOINT [ "./secretgw" ]
