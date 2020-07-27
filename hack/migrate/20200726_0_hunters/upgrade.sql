CREATE TABLE hunters (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT hunters_id_pk
            PRIMARY KEY,
    user_id bigint NOT NULL,
    name text  NOT NULL,
    skills text[] NOT NULL,
    about text NOT NULL,
    description_addr text NOT NULL,
    email text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX hunters_user_id_uindex ON hunters (user_id);