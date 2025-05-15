-- Create boards table
CREATE TABLE boards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

-- Create feedback table
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    board_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category_id INT NOT NULL,
    upvotes INT DEFAULT 0,
    downvotes INT DEFAULT 0,
    FOREIGN KEY (board_id) REFERENCES boards(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Create votes table
CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    feedback_id INT NOT NULL,
    user_id INT NOT NULL,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('upvote', 'downvote')),
    FOREIGN KEY (feedback_id) REFERENCES feedback(id),
    UNIQUE (feedback_id, user_id) -- Ensure a user can only vote once per feedback
);
