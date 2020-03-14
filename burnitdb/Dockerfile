FROM golang:alpine AS builder
# Create a user.
RUN adduser -D -g '' burnitdb

# Add compilation further on...

FROM scratch
# Set default port.
ARG VERSION
ARG PORT=3001
# Get user settings.
COPY --from=builder /etc/passwd /etc/passwd

COPY release/bin/burnitdb /usr/bin/burnitdb

EXPOSE ${PORT}

USER burnitdb
ENTRYPOINT [ "/usr/bin/burnitdb" ]
