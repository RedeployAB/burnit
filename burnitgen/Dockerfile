FROM golang:alpine AS builder
# Create a user.
RUN adduser -D -g '' burnitgen

# Add compilation further on...

FROM scratch
# Set default port.
ARG VERSION
ARG PORT=3002
# Get user settings.
COPY --from=builder /etc/passwd /etc/passwd

COPY release/bin/burnitgen /usr/bin/burnitgen

EXPOSE ${PORT}

USER burnitgen
ENTRYPOINT [ "/usr/bin/burnitgen" ]
