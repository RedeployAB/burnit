FROM golang@sha256:4d8abd16b03209b30b48f69a2e10347aacf7ce65d8f9f685e8c3e20a512234d9 AS builder
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
