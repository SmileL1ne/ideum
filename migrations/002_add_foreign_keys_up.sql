/*
    TODO:
    - Redo eveything here, add and delete necessary foreign keys
*/

-- Add foreign key to 'posts' table
ALTER TABLE posts ADD COLUMN user_id INTEGER REFERENCES users(id);

-- Add foreign keys to 'post_likes' table
ALTER TABLE post_likes ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE post_likes ADD COLUMN post_id INTEGER REFERENCES posts(id);

-- Add foreign keys to 'comments' table
ALTER TABLE comments ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE comments ADD COLUMN post_id INTEGER REFERENCES posts(id);

-- Add foreign keys to 'comment_likes' table
ALTER TABLE comment_likes ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE comment_likes ADD COLUMN comment_id INTEGER REFERENCES comments(id);
