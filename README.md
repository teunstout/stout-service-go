# Welcome to the documentation

This documentation is personal and i document some things in here that i might need or things that aren't often used. So this is basically treated as a cheat sheet ;

## Content

- [Welcome to the documentation](#welcome-to-the-documentation)
  - [Content](#content)
  - [Architecture](#architecture)
  - [⚙️ Commands](#️-commands)
    - [Quickstart Docker Compose](#quickstart-docker-compose)
    - [Generate new module](#generate-new-module)
  - [🔗 Sources](#-sources)

## Architecture

```mermaid
flowchart LR
    user(["👤 User"])

    subgraph public["Public Edge"]
        frontend["Frontend"]
        idpService["IDP Service"]
        databaseIdp[("Postgres")]
    end

    subgraph internal["Internal Network"]
        nginx["Nginx<br/><i>Reverse Proxy</i>"]
        contentService["Content Service"]
        jishoService["Jisho Service"]
        database[("Postgres")]
    end

    user --> frontend
    frontend --> idpService
    frontend --> nginx
    idpService --> databaseIdp
    nginx --> idpService
    nginx --> contentService
    nginx --> jishoService
    contentService --> database
    jishoService --> database

    classDef userStyle fill:#e1f5ff,stroke:#0288d1,stroke-width:2px,color:#01579b
    classDef edgeStyle fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100
    classDef internalStyle fill:#f1f8e9,stroke:#558b2f,stroke-width:2px,color:#33691e
    classDef proxyStyle fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f
    classDef dbStyle fill:#ede7f6,stroke:#5e35b1,stroke-width:2px,color:#311b92

    class user userStyle
    class frontend,idpService edgeStyle
    class contentService,jishoService internalStyle
    class nginx proxyStyle
    class database dbStyle
    class databaseIdp dbStyle
```

## ⚙️ Commands

### Quickstart Docker Compose

```bash
docker compose up
```

### Generate new module

```bash
go mod init example/package
```

## 🔗 Sources

|                             |                                                                    |
| --------------------------- | ------------------------------------------------------------------ |
| Name                        | Url                                                                |
| net/http documentation      | <https://go.dev/doc/articles/wiki/>                                |
| OAuth2 Client side          | <https://pkg.go.dev/golang.org/x/oauth2#section-readme>            |
| OAuth2 Client server side   | <https://pkg.go.dev/github.com/go-oauth2/oauth2/v4#section-readme> |
| OAuth2 Client with Postgres | <https://github.com/vgarvardt/go-oauth2-pg>                        |
