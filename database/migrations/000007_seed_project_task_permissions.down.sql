DELETE
FROM permissions
WHERE key IN (
              'projects.create', 'projects.update', 'projects.delete',
              'tasks.create', 'tasks.update', 'tasks.delete', 'tasks.assign'
    );
