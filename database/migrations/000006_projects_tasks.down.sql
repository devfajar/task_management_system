DROP INDEX IF EXISTS idx_tasks_status_due;
DROP INDEX IF EXISTS idx_tasks_assignee;
DROP INDEX IF EXISTS idx_tasks_project;
DROP INDEX IF EXISTS idx_projects_owner;

DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
