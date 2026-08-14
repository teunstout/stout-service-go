# 🚀 nginx

Reverse proxy for stout-service-go. This is the only container in `docker-compose.yaml` that publishes a host port (`8080`) — every other service is only reachable from inside the compose network, through here.

## Routing

Requests are routed by path prefix; the prefix is stripped before proxying, and CORS headers (`cors.conf`, wide-open `*`) are added on every proxied route. Anything that doesn't match falls through to a plain `404`.

| Path prefix | Proxied to |
| --- | --- |
| `/idp-service/*` | `idp_service:8080/*` (see [idp-service](../idp-service/README.MD)) |
| `/content-service/*` | `content_service:8080/*` (see [content-service](../content-service/README.md)) |
| `/jisho-service/*` | `jisho_service:8080/*` (see [jisho-service](../jisho-service/README.MD)) |
| `/health` | always returns `200` (this only checks nginx itself, not the upstream services) |

For example, `GET /idp-service/v1/login` is proxied to `idp_service:8080/v1/login`.

**Health:**

```shell
curl -i -X GET http://localhost:8080/health
```

**Example proxied requests** (see each service's README for the full curl list):

```shell
curl -i -X GET http://localhost:8080/idp-service/.well-known/jwks.json

curl -i -X GET http://localhost:8080/content-service/v1/content/cv --output cv.pdf

curl -i -X GET "http://localhost:8080/jisho-service/v1/search?keyword=犬" \
  -H "Authorization: Bearer <jwt-from-idp-login>"
```

## Config files

- `nginx.conf` — server block with the routing table above.
- `cors.conf` — included by every proxied `location`; sets wide-open CORS headers and handles `OPTIONS` preflight requests.
- `Dockerfile` — `FROM nginx:stable-bookworm`, just copies in the two config files above.

## Local edits

After changing `nginx.conf` or `cors.conf`, rebuild and restart the container to pick them up:

```shell
docker compose up -d --build nginx
```
