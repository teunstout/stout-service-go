-- idp-server.
CREATE TABLE
    IF NOT EXISTS accounts (
        id SERIAL PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

CREATE TABLE
    IF NOT EXISTS session_tokens (
        id SERIAL PRIMARY KEY,
        token VARCHAR(64) NOT NULL,
        account_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        expires_at TIMESTAMP NOT NULL,
        FOREIGN KEY (account_id) REFERENCES accounts (id),
        UNIQUE (token, account_id)
    );

CREATE INDEX IF NOT EXISTS idx_session_account_id ON session_tokens (account_id);

CREATE INDEX IF NOT EXISTS idx_session_token ON session_tokens (token);

CREATE TABLE
    IF NOT EXISTS csrf_tokens (
        id SERIAL PRIMARY KEY,
        token VARCHAR(64) NOT NULL,
        account_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        expires_at TIMESTAMP NOT NULL,
        FOREIGN KEY (account_id) REFERENCES accounts (id),
        UNIQUE (token, account_id)
    );

CREATE INDEX IF NOT EXISTS idx_csrf_account_id ON csrf_tokens (account_id);

CREATE INDEX IF NOT EXISTS idx_csrf_token ON csrf_tokens (token);

-- jisho-service.
CREATE TABLE
    IF NOT EXISTS jisho (id SERIAL PRIMARY KEY);

CREATE TABLE
    IF NOT EXISTS member_search_history (
        id SERIAL PRIMARY KEY,
        member_id INTEGER NOT NULL,
        search_term TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );


CREATE INDEX IF NOT EXISTS idx_member_search_history_member_id ON member_search_history (member_id);
CREATE INDEX IF NOT EXISTS idx_member_search_history_search_term ON member_search_history (search_term);
