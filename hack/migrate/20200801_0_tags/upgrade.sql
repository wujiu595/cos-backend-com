CREATE TABLE tags
(
    id         BIGINT                   DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT tags_id_pk
            PRIMARY KEY,
    name       TEXT                                                     NOT NULL,
    source     TEXT                                                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    deleted    BOOLEAN                  DEFAULT FALSE                   NOT NULL
);

COMMENT ON COLUMN tags.source IS 'skills';

CREATE UNIQUE INDEX tags_source_name
    ON tags (source,name);