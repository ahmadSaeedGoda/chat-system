CREATE TABLE IF NOT EXISTS chat.messages (
    user TEXT,
    timestamp TIMESTAMP,
    id UUID,
    sender TEXT,
    recipient TEXT,
    content TEXT,
    PRIMARY KEY (user, timestamp, id)
) WITH CLUSTERING ORDER BY (timestamp DESC);
