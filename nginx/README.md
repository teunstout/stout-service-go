# 🚀 proxy

Proxy for the services from teun stout

**Health:**

```shell
curl --request GET --url http://ah-graphql-router-proxy:4000/health
```

## Filtering requests

```shell
# Build with lua
docker build --build-arg ENABLED_MODULES="ndk lua" -t nginx .
```