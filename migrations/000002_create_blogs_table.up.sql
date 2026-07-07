-- blogs up

CREATE TABLE IF NOT EXISTS blogs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    category    VARCHAR(100),
    tags        TEXT[] NOT NULL DEFAULT '{}',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common access patterns.
CREATE INDEX idx_blogs_user_id    ON blogs (user_id);
CREATE INDEX idx_blogs_updated_at ON blogs (updated_at DESC);

