-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user' NOT NULL;

-- Create board_members table to track user roles within boards
CREATE TABLE board_members (
    id SERIAL PRIMARY KEY,
    board_id INT NOT NULL,
    user_id INT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- 'stakeholder' or 'user'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(board_id, user_id)
);

-- Create index for faster lookups
CREATE INDEX idx_board_members_board_id ON board_members(board_id);
CREATE INDEX idx_board_members_user_id ON board_members(user_id);

-- Update the first user as app_admin (You'll need to run this manually once you have a user)
-- UPDATE users SET role = 'app_admin' WHERE id = 1;
