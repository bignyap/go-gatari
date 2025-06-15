-- name: ListOrganization :many
SELECT 
    organization.*, organization_type.organization_type_name, 
    COUNT(organization_id) OVER() AS total_items 
FROM organization
INNER JOIN organization_type ON organization.organization_type_id = organization_type.organization_type_id
WHERE (? IS NULL OR organization.organization_id = ?)
ORDER BY organization_id DESC
LIMIT ? OFFSET ?;

-- name: CreateOrganization :execresult 
INSERT INTO organization (
    organization_name, organization_created_at, organization_updated_at, 
    organization_realm, organization_country, organization_support_email,
    organization_active, organization_report_q, organization_config,
    organization_type_id
) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateOrganizations :copyfrom
INSERT INTO organization (
    organization_name, organization_created_at, organization_updated_at, 
    organization_realm, organization_country, organization_support_email,
    organization_active, organization_report_q, organization_config,
    organization_type_id
) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateOrganization :execresult
UPDATE organization
SET 
    organization_name = ?,
    organization_updated_at = ?,
    organization_realm = ?,
    organization_country = ?,
    organization_support_email = ?,
    organization_active = ?,
    organization_report_q = ?,
    organization_config = ?,
    organization_type_id = ?
WHERE organization_id = ?;

-- name: DeleteOrganizationById :exec
DELETE FROM organization
WHERE organization_id = ?;