CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE posts (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(150) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) 
);

CREATE TABLE post_likes (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    is_like BOOLEAN NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE comments (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE comment_likes (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    is_like BOOLEAN NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE categories (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    created DATETIME NOT NULL,
);

-- Bridge table, connects posts and categories tables
CREATE TABLE posts_categories (
    post_id INTEGER,
    category_id INTEGER,

    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

