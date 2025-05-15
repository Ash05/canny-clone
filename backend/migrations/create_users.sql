-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    picture TEXT,
    provider VARCHAR(50) NOT NULL, -- e.g., "google", "github"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Add foreign key constraints to related tables
ALTER TABLE comments ADD CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE comment_replies ADD CONSTRAINT fk_comment_replies_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE comment_likes ADD CONSTRAINT fk_comment_likes_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE votes ADD CONSTRAINT fk_votes_user FOREIGN KEY (user_id) REFERENCES users(id);
