CREATE TABLE stats_db.views (
    post_id UInt64,
    appraiser_id UInt64,
) ENGINE = MergeTree()
ORDER BY post_id
PRIMARY KEY post_id;

CREATE TABLE stats_db.likes (
    post_id UInt64,
    appraiser_id UInt64,
) ENGINE = MergeTree()
ORDER BY post_id
PRIMARY KEY post_id;
