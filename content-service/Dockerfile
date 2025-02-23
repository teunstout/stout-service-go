FROM debian:bookworm

WORKDIR /app

COPY --chown=1000:1000 ./content /app/content
COPY --chown=1000:1000 ./assets /app/assets/
USER 1000:1000

CMD ["./content"]
