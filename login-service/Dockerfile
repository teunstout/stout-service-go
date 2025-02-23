FROM debian:bookworm

WORKDIR /app

COPY --chown=1000:1000 ./login /app/login

USER 1000:1000
CMD ["./login"]
