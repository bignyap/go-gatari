-- name: ListOrganization :many
SELECT 
    organization.*, 
    organization_type.organization_type_name, 
    COUNT(*) OVER() AS total_items
FROM organization
INNER JOIN organization_type 
    ON organization.organization_type_id = organization_type.organization_type_id
WHERE ($1::INTEGER = 0 OR organization.organization_id = $2)
ORDER BY organization.organization_id DESC
LIMIT $3 OFFSET $4;

-- name: CreateOrganization :one 
INSERT INTO organization (
    organization_name, organization_created_at, organization_updated_at, 
    organization_realm, organization_country, organization_support_email,
    organization_active, organization_report_q, organization_config,
    organization_type_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING organization_id;

-- name: CreateOrganizations :copyfrom
INSERT INTO organization (
    organization_name, organization_created_at, organization_updated_at, 
    organization_realm, organization_country, organization_support_email,
    organization_active, organization_report_q, organization_config,
    organization_type_id
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdateOrganization :execresult
UPDATE organization
SET 
    organization_name = $1,
    organization_updated_at = $2,
    organization_realm = $3,
    organization_country = $4,
    organization_support_email = $5,
    organization_active = $6,
    organization_report_q = $7,
    organization_config = $8,
    organization_type_id = $9
WHERE organization_id = $10;

-- name: DeleteOrganizationById :exec
DELETE FROM organization
WHERE organization_id = $1;