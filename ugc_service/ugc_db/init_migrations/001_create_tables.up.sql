SET TIMEZONE="Europe/Moscow";

CREATE TABLE posts (
    post_id BIGSERIAL NOT NULL,
    author_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    create_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    update_timestamp TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (post_id)
);

