-- admin: all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- moderator: read any user + update any user + read roles/permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
  ON p.name IN ('users.read.any', 'users.update.any', 'roles.read', 'permissions.read')
WHERE r.name = 'moderator'
ON CONFLICT DO NOTHING;

-- user: read/update/delete own record only
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.name IN ('users.read.self', 'users.update.self', 'users.delete.self')
WHERE r.name = 'user'
ON CONFLICT DO NOTHING;
