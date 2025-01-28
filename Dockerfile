# golang:1.23.4-alpine3.21 SHA digest.
FROM --platform=$BUILDPLATFORM golang@sha256:6c5c9590f169f77c8046e45c611d3b28fe477789acd8d3762d23d4744de69812 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BIN=burnit

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ENV USER=${BIN}
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nohome" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"


FROM scratch

ARG TARGETOS
ARG TARGETARCH
ARG BIN=burnit
ARG PORT=3000

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY ./build/$TARGETOS/$TARGETARCH/${BIN} /${BIN}

EXPOSE ${PORT}

USER ${BIN}:${BIN}

ENTRYPOINT [ "/burnit" ]
