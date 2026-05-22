INSERT INTO permissions (name, description) VALUES
    ('users.read.self',   'Read own user record'),
    ('users.read.any',    'Read any user record'),
    ('users.create',      'Create user records'),
    ('users.update.self', 'Update own user record'),
    ('users.update.any',  'Update any user record'),
    ('users.delete.self', 'Delete own user record'),
    ('users.delete.any',  'Delete any user record'),
    ('roles.read',        'Read roles'),
    ('roles.write',       'Create or update roles'),
    ('roles.delete',      'Delete roles'),
    ('permissions.read',  'Read permissions'),
    ('permissions.write', 'Create or update permissions'),
    ('permissions.delete','Delete permissions')
ON CONFLICT (name) DO NOTHING;
