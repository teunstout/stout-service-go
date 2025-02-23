FROM debian:bookworm

WORKDIR /app

COPY --chown=1000:1000 ./authorization /app/authorization

USER 1000:1000
CMD ["./authorization"]
