CREATE TABLE users
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT users_id_pk
            PRIMARY KEY,
    name       TEXT                                                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    deleted    BOOLEAN                  DEFAULT FALSE                   NOT NULL
);
CREATE UNIQUE INDEX users_name ON users (name);

CREATE TABLE roles
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT roles_id_pk
            PRIMARY KEY,
    name       TEXT                                                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    deleted    BOOLEAN                  DEFAULT FALSE                   NOT NULL
);
CREATE UNIQUE INDEX roles_name ON roles (name);

CREATE TABLE users_roles_rel
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT users_roles_rel_pk
            PRIMARY KEY,
    role_id    BIGINT                                                   NOT NULL,
    user_id    BIGINT                                                   NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL
);

CREATE UNIQUE INDEX users_roles_rel_role_id_user_id ON roles (created_at,updated_at);