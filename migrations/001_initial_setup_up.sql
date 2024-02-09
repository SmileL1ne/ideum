/*
    TODO: 
    - Create indexes for each table
    - Redo every table with 'REWORK' tag (see below) by blueprint in schema (lucidchart)  
*/

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
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    created DATETIME NOT NULL
);

-- REWORK
-- Bridge table, connects posts and categories tables
CREATE TABLE IF NOT EXISTS posts_categories (
    post_id INTEGER REFERENCES posts(id),
    category_id INTEGER REFERENCES categories(id),

    PRIMARY KEY (post_id, category_id)
);

-- REWORK
-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT NOT NULL PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);