FROM golang:alpine AS builder
# Create a user.
RUN adduser -D -g '' app

# Add compilation further on...

FROM scratch
# Set default port.
ARG VERSION
ARG PORT=3002
# Get user settings.
COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /app
COPY release/bin/secretgen /app/

EXPOSE ${PORT}

USER app
ENTRYPOINT [ "./secretgen" ]
