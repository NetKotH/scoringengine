ARG PACKAGE=github.com/netkoth/scoringengine
FROM brimstone/golang-musl as builder

FROM scratch
EXPOSE 80 443 53/udp 67/udp
ENTRYPOINT ["/scoringengine"]
COPY --from=builder /app /scoringengine
