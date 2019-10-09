FROM scratch

ARG VERSION
ARG PORT=3001

WORKDIR /app

COPY release/bin/secretdb /app/

EXPOSE ${PORT}

ENTRYPOINT [ "./secretdb" ]
