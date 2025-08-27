-- seed minimal, tanpa menyebut ID (biar default v4 yang mengisi)
INSERT INTO permissions (key, description)
VALUES ('users.force_delete', 'Hard delete user (admin only)')
ON CONFLICT DO NOTHING;
