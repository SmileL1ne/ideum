/*
    TODO: 
    - Create indexes for each table
    - Redo every table with 'REWORK' tag (see below) by blueprint in schema (lucidchart)  
*/

-- Turn on foreign keys (they are disabled by default for backward compatibility)
PRAGMA foreign_keys = ON;

-- REWORK
CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(150) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
CREATE TABLE IF NOT EXISTS post_reactions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    is_like BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
CREATE TABLE IF NOT EXISTS comment_reactions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    is_like BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL
);

-- REWORK
-- Bridge table, connects posts and tags tables
CREATE TABLE IF NOT EXISTS posts_tags (
    post_id INTEGER REFERENCES posts(id),
    tag_id INTEGER REFERENCES tags(id),
    created_at DATETIME NOT NULL,

    PRIMARY KEY (post_id, tag_id)
);

-- REWORK
-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry DATETIME NOT NULL
);