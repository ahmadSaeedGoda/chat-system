CREATE TABLE IF NOT EXISTS chat.messages_by_recipient (
    recipient TEXT,
    id UUID,
    sender TEXT,
    timestamp TIMESTAMP,
    content TEXT,
    PRIMARY KEY (recipient, timestamp, id)
);
