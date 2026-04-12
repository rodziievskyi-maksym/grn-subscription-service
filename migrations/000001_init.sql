CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    repository VARCHAR(255) NOT NULL,
    last_seen_tag VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    UNIQUE(email, repository)
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_repository ON subscriptions(repository);
CREATE INDEX IF NOT EXISTS idx_subscriptions_email ON subscriptions(email);