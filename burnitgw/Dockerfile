FROM golang:alpine AS builder
# Create a user.
RUN adduser -D -g '' burnitgw

# Add compilation further on...

FROM scratch
# Set default port.
ARG VERSION
ARG PORT=3000
# Get user settings.
COPY --from=builder /etc/passwd /etc/passwd

COPY release/bin/burnitgw /usr/bin/burnitgw

EXPOSE ${PORT}

USER burnitgw
ENTRYPOINT [ "/usr/bin/burnitgw" ]
