CREATE TABLE IF NOT EXISTS chat.messages_by_sender (
    sender TEXT,
    id UUID,
    recipient TEXT,
    timestamp TIMESTAMP,
    content TEXT,
    PRIMARY KEY (sender, timestamp, id)
);
