-- Boards & Columns permissions
INSERT INTO permissions (key, description)
VALUES ('boards.create', 'Create board'),
       ('boards.update', 'Update board'),
       ('boards.delete', 'Delete board'),
       ('columns.create', 'Create column'),
       ('columns.update', 'Update column'),
       ('columns.delete', 'Delete column')
ON CONFLICT DO NOTHING;
