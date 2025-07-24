-- Add configuration permissions
INSERT INTO permissions (name, description, resource, action, status_id, created_at, updated_at)
VALUES
('Create Configurations', 'Permission to create new configurations', 'configurations', 'create', 0, NOW(), NOW()),
('Read Configurations', 'Permission to read configurations', 'configurations', 'read', 0, NOW(), NOW()),
('Update Configurations', 'Permission to update configurations', 'configurations', 'update', 0, NOW(), NOW()),
('Delete Configurations', 'Permission to delete configurations', 'configurations', 'delete', 0, NOW(), NOW());

-- Assign permissions to Editor group (assuming group_id = 2)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 2, id FROM permissions WHERE resource = 'configurations' AND action IN ('create', 'read', 'update');

-- Assign read permission to Viewer group (assuming group_id = 3)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 3, id FROM permissions WHERE resource = 'configurations' AND action = 'read';

-- Assign read permission to General client group (assuming group_id = 4)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 4, id FROM permissions WHERE resource = 'configurations' AND action = 'read';

-- Assign all permissions to Admin group (assuming group_id = 1)
INSERT INTO group_permissions (group_id, permission_id)
SELECT 1, id FROM permissions WHERE resource = 'configurations';
