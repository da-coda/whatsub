CREATE TABLE subreddits (
    sub_name text PRIMARY KEY,
    sub_full_name text,
    subscriber_count int,
    created_date date
);

CREATE INDEX subreddits_subscriber_count ON subreddits (subscriber_count);

CREATE TABLE temp (sub_name text, created_date text, subs text);

COPY temp FROM '/subreddits_2021-02-19.csv' WITH (FORMAT CSV);

INSERT INTO subreddits
SELECT temp.sub_name, 'r/' || temp.sub_name, temp.subs::int, to_date(temp.created_date, 'YYYY-MM-DD') FROM temp WHERE sub_name IS NOT NULL AND subs IS NOT NULL AND created_date IS NOT NULL;

DROP TABLE temp;