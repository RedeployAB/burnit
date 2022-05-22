FROM golang@sha256:42d35674864fbb577594b60b84ddfba1be52b4d4298c961b46ba95e9fb4712e8 AS builder

ARG OS=linux
ARG ARCH=amd64
ARG BIN=burnitgw

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ENV USER=${BIN}
ENV UID=10001

# Create a user.
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nohome" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src/${BIN}
COPY . .

RUN mv build/common ../

RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -o build/${BIN} -ldflags="-w -s" -trimpath cmd/${BIN}/main.go

FROM scratch

ARG BIN=burnitgw
ARG PORT=3000

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /src/${BIN}/build/${BIN} /usr/bin/${BIN}

EXPOSE ${PORT}

USER ${BIN}:${BIN}

ENTRYPOINT [ "/usr/bin/burnitgw" ]
