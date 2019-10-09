FROM scratch

ARG VERSION
ARG PORT=3002

WORKDIR /app

COPY release/bin/secretgen /app/

EXPOSE ${PORT}

ENTRYPOINT [ "./secretgen" ]
