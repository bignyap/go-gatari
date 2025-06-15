-- name: ListSubscriptionTier :many
SELECT *, COUNT(subscription_tier_id) OVER() AS total_items  
FROM subscription_tier
WHERE tier_archived = ?
ORDER BY subscription_tier_id DESC
LIMIT ? OFFSET ?;

-- name: ArchiveExistingSubscriptionTier :exec
UPDATE subscription_tier
SET tier_archived = TRUE
WHERE tier_name = ?;

-- name: CreateSubscriptionTier :execresult 
INSERT INTO subscription_tier (tier_name, tier_description, tier_created_at, tier_updated_at)
VALUES (?, ?, ?, ?);

-- name: CreateSubscriptionTiers :copyfrom
INSERT INTO subscription_tier (tier_name, tier_description, tier_created_at, tier_updated_at) 
VALUES (?, ?, ?, ?);

-- name: DeleteSubscriptionTierById :exec
DELETE FROM subscription_tier
WHERE subscription_tier_id = ?;