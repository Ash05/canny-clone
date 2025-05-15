-- Create comments table
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    feedback_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    FOREIGN KEY (feedback_id) REFERENCES feedback(id)
);

-- Create comment_replies table
CREATE TABLE comment_replies (
    id SERIAL PRIMARY KEY,
    comment_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);

-- Create comment_likes table to track who liked a comment
CREATE TABLE comment_likes (
    id SERIAL PRIMARY KEY,
    comment_id INT,
    reply_id INT,
    user_id INT NOT NULL,
    is_like BOOLEAN NOT NULL, -- true for like, false for dislike
    CHECK ((comment_id IS NULL) != (reply_id IS NULL)), -- Either comment_id or reply_id must be non-null, but not both
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    FOREIGN KEY (reply_id) REFERENCES comment_replies(id),
    UNIQUE (comment_id, user_id), -- Ensure a user can only have one reaction per comment
    UNIQUE (reply_id, user_id)    -- Ensure a user can only have one reaction per reply
);