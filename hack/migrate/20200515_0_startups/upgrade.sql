SET SEARCH_PATH = comunion;

CREATE TABLE users
(
    id          BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT users_id_pk
            PRIMARY KEY,
    wallet_addr TEXT                                                     NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    is_hunter   BOOL                     DEFAULT FALSE                   NOT NULL
);

CREATE TABLE start_ups
(
    id               BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT start_ups_id_pk
            PRIMARY KEY,
    name             TEXT                                                     NOT NULL,
    mission          TEXT,
    logo             TEXT                                                     NOT NULL,
    tx_id            TEXT                                                     NOT NULL,
    blockNum         BIGINT,
    description_addr TEXT                                                     NOT NULL,
    category_id      BIGINT                                                   NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    state            INT                      DEFAULT 0                       NOT NULL,
    isIRO            BOOL                     DEFAULT FALSE                   NOT NULL
);

CREATE UNIQUE INDEX start_ups_tx_id ON start_ups (tx_id);
COMMENT ON COLUMN start_ups.state IS '0 创建中,1 已创建,2 未确认到tx产生,3 上链失败，4 已设置';

CREATE TABLE categories
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT categories_id_pk
            PRIMARY KEY,
    name       TEXT                                                     NOT NULL,
    code       TEXT                                                     NOT NULL,
    source     TEXT                                                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    deleted    BOOL                     DEFAULT FALSE                   NOT NULL
);

CREATE UNIQUE INDEX categories_name ON categories (name);
CREATE UNIQUE INDEX categories_code ON categories (code);

COMMENT ON COLUMN categories.source IS 'start_up';

CREATE TABLE access_keys
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT access_keys_id_pk
            PRIMARY KEY,
    key        TEXT                                                     NOT NULL,
    secret     TEXT                                                     NOT NULL,
    uid        BIGINT                                                   NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    state      SMALLINT                 DEFAULT 0                       NOT NULL
);

CREATE UNIQUE INDEX access_keys_key ON access_keys (id);