FROM caddy:2.6.4-builder AS builder

RUN xcaddy build fix-empty-copy-headers

FROM caddy:2.7.5

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
