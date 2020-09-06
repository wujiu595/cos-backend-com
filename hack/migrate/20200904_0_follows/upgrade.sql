CREATE TABLE startups_follows_rel
(
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT startups_follows_rel_pk
            PRIMARY KEY,
    startup_id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX startups_follows_rel_startup_id_user_id ON startups_follows_rel(startup_id,user_id);