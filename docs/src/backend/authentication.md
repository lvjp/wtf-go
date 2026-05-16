# Authentication

The API uses **Bearer token** authentication. Tokens are unique opaque token IDs managed entirely server-side and stored via an external store.

> **Note:** The current implementation uses `NewMemoryStore`, which is suitable for development and testing only. All tokens are lost on server restart.

## Creating a token

```http
POST /api/v0/auth/token
Content-Type: application/json

{
  "subject": "my-api-client"
}
```

The request body is optional. `subject` is an arbitrary label identifying the token holder (default: empty string).

The token lifetime is set server-side via the `auth.token_ttl` configuration key (default: `24h`).

**Response `201`:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "not_after": "2026-05-15T12:00:00Z"
}
```

Use the `id` value as the Bearer token on protected endpoints.

## Using a token

Add an `Authorization` header to every request targeting a protected route:

```http
GET /api/v0/some/protected/resource
Authorization: Bearer 550e8400-e29b-41d4-a716-446655440000
```

A missing, malformed, or expired token returns `401 Unauthorized`.

A token is considered valid **at** its `not_after` instant, and expires only **strictly after** it.

## Revoking a token

Revokes the token used for authentication. This endpoint requires authentication.

```http
DELETE /api/v0/auth/token
Authorization: Bearer 550e8400-e29b-41d4-a716-446655440000
```

Returns `204 No Content` on success, `401` if the token is missing, `403` if the token is invalid.
