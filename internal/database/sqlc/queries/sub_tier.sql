-- name: ListSubscriptionTier :many
SELECT *, COUNT(subscription_tier_id) OVER() AS total_items  
FROM subscription_tier
WHERE tier_archived = $1
ORDER BY subscription_tier_id DESC
LIMIT $2 OFFSET $3;

-- name: ArchiveExistingSubscriptionTier :exec
UPDATE subscription_tier
SET tier_archived = TRUE
WHERE tier_name = $1;

-- name: CreateSubscriptionTier :one 
INSERT INTO subscription_tier (
    tier_name, tier_description, tier_created_at, tier_updated_at
)
VALUES ($1, $2, $3, $4)
RETURNING subscription_tier_id;

-- name: CreateSubscriptionTiers :copyfrom
INSERT INTO subscription_tier (
    tier_name, tier_description, tier_created_at, tier_updated_at
) 
VALUES ($1, $2, $3, $4);

-- name: DeleteSubscriptionTierById :exec
DELETE FROM subscription_tier
WHERE subscription_tier_id = $1;