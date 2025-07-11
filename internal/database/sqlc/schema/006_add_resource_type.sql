-- +goose Up
CREATE TABLE permission_type (
  permission_code VARCHAR(10) PRIMARY KEY,           -- 3-letter code, e.g. 'CRT', 'RED'
  permission_name VARCHAR(50) NOT NULL,
  permission_description TEXT
);

-- Insert standard and extended permissions using 3-letter codes
INSERT INTO permission_type (permission_code, permission_name, permission_description) VALUES
  ('CRT', 'Create',     'Permission to create new records'),
  ('RED', 'Read',       'Permission to read or view records'),
  ('UPD', 'Update',     'Permission to update or modify records'),
  ('DEL', 'Delete',     'Permission to delete records'),
  ('LST', 'List',       'Permission to list all records'),
  ('EXE', 'Execute',    'Permission to execute actions like report generation'),
  ('UPL', 'Upload',     'Permission to upload files or data'),
  ('DNL', 'Download',   'Permission to download files or data'),
  ('CFG', 'Configure',  'Permission to change configuration settings'),
  ('ADM', 'Administer', 'Permission to manage users or settings'),
  ('ANL', 'Analyze',    'Permission to run analysis operations'),
  ('SYN', 'Sync',       'Permission to synchronize data'),
  ('ARC', 'Archive',    'Permission to archive resources'),
  ('RES', 'Restore',    'Permission to restore archived or deleted resources'),
  ('SHR', 'Share',      'Permission to share resource access');

-- +goose Down
-- Then drop the permission_type table
DROP TABLE IF EXISTS permission_type;
