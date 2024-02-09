/*
    TODO:
    - Redo eveything here, add and delete necessary foreign keys
*/

-- Add foreign key to 'posts' table
ALTER TABLE posts ADD COLUMN user_id INTEGER REFERENCES users(id);

-- Add foreign keys to 'post_reactions' table
ALTER TABLE post_reactions ADD COLUMN post_id INTEGER REFERENCES posts(id);
ALTER TABLE post_reactions ADD COLUMN user_id INTEGER REFERENCES users(id);

-- Add foreign keys to 'comments' table
ALTER TABLE comments ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE comments ADD COLUMN post_id INTEGER REFERENCES posts(id);

-- Add foreign keys to 'comment_likes' table
ALTER TABLE comment_reactions ADD COLUMN comment_id INTEGER REFERENCES comments(id);
ALTER TABLE comment_reactions ADD COLUMN user_id INTEGER REFERENCES users(id);
