DELETE
FROM permissions
WHERE key IN (
              'boards.create', 'boards.update', 'boards.delete',
              'columns.create', 'columns.update', 'columns.delete'
    );
