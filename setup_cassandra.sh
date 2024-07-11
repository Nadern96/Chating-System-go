#!/bin/bash

# Check if port argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <cassandra_port>"
    exit 1
fi

# Assign the port argument to a variable
port=$1

# Command to connect to Cassandra using cqlsh with specified port
cqlsh localhost $port <<EOF

-- Create keyspace 'auth'
CREATE KEYSPACE IF NOT EXISTS auth WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

-- Switch to keyspace 'auth'
USE auth;

-- DROP TABLE 
DROP TABLE IF EXISTS user;

-- Create table 'user' in keyspace 'auth'
CREATE TABLE IF NOT EXISTS user (
    username text,
    id uuid,
    createdat timestamp,
    password text,
    email text,
    PRIMARY KEY (id)
);

-- Create an index on 'email' column
CREATE INDEX IF NOT EXISTS ON user (email);


-- Create keyspace 'chat'
CREATE KEYSPACE IF NOT EXISTS chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

-- Switch to keyspace 'chat'
USE chat;

-- DROP TABLE 
DROP TABLE IF EXISTS chat;

-- Create table 'chat' in keyspace 'chat'
CREATE TABLE IF NOT EXISTS chat (
    chatId UUID,
    fromUserId UUID,
    toUserId UUID,
    PRIMARY KEY (chatId)
);

CREATE INDEX IF NOT EXISTS ON chat (fromUserId);
CREATE INDEX IF NOT EXISTS ON chat (toUserId);


-- DROP TABLE 
DROP TABLE IF EXISTS message;

-- Create table 'message' in keyspace 'chat'
CREATE TABLE IF NOT EXISTS message (
    chatId UUID,
    messageId TIMEUUID,
    fromUserId UUID,
    toUserId UUID,
    content TEXT,
    createdAt TIMESTAMP,
    PRIMARY KEY (chatId, messageId)
);

EOF
