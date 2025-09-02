-- permissions dasar untuk projects & tasks
INSERT INTO permissions (key, description)
VALUES ('projects.create', 'Create project'),
       ('projects.update', 'Update project'),
       ('projects.delete', 'Delete project'),
       ('tasks.create', 'Create task'),
       ('tasks.update', 'Update task'),
       ('tasks.delete', 'Delete task'),
       ('tasks.assign', 'Assign task to user')
ON CONFLICT DO NOTHING;
