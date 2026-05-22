INSERT INTO roles (name, description) VALUES
    ('admin',     'Full system access'),
    ('moderator', 'Content moderation access'),
    ('user',      'Standard user access')
ON CONFLICT (name) DO NOTHING;
