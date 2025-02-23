FROM debian:bookworm

WORKDIR /app

COPY --chown=1000:1000 ./jisho /app/jisho
USER 1000:1000

CMD ["./jisho"]
