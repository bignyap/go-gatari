# Gatekeeper Functionality

## Request: `/validateRequest`

### Inputs:
1. **OrgName**
2. **Method**
3. **Path**

### Goal:
1. **Get Org Details**: Verifies if the organization exists and applies overall rate limiting.
2. **Get Subscription Details**: Checks if the subscription is valid.
3. **Get Usage**: Confirms if the organization has credits left.
4. **Get Endpoint Details**: Confirms rate limiting for the API.

---

## Response: `/recordUsage`

### Inputs:
1. **OrgName**
2. **Method**
3. **Path**

### Goal:
1. Retrieves organization, subscription, and endpoint details to update usage/credits.

---

## Caching

### Types:
1. **Get Endpoint, Subscription, and Organization**: Read-only with occasional writes.
2. **Endpoint Loading**: Read with occasional writes.
3. **Update API Count and Credit Usage**:
	- **Total API Count**: Read + Write (used for rate limiting at the organization level).
	- **Total Credit Usage**: Read + Write (used to restrict credit usage).
	- **API-wise Count**: Read + Write (used for rate limiting at the API level).

---

## Caching Strategy

### Read-only + Occasional Write:
- **Pub/Sub Strategy**

### Frequent Write:
1. Read from the database.
2. Store in Local + Redis.
3. Write to Local.
4. Write to Redis (5-second aggregate).
5. Write to Database (30-second aggregate).
6. Invalidate Local and Redis.
7. Repeat.