# User API

## Create a user

`POST /v1/user` accepts a JSON user object. `name`, `email`, and `password` are used to create the account. Registrations default to the `customer` role; any other client-supplied role is ignored.

The optional `X-Request-Origin` header supports these values:

| Value | Behavior |
| --- | --- |
| Omitted | Preserve the existing `finance-api` and `akademia-api` `user.created` publications. |
| `barber-app` | Preserve both existing publications and also publish `user.created` to the `barber-app` topic. |

Other non-empty values receive `400 Bad Request`. To request the `barber` role, send `"role": "barber"` with a valid admin JWT in the `Authorization: Bearer <token>` header. The API verifies the token signature with its configured RSA public key and checks the signed `role` claim. A request body or header that claims to be admin is not authorization. Requests without valid admin authorization receive `403 Forbidden` for a barber role request.

A successful creation returns `201 Created` with the created ID:

```json
{
  "message": "User created successfully",
  "user_id": "<created-user-id>",
  "user": {
    "id": "<created-user-id>",
    "name": "<user-name>",
    "email": "<user-email>",
    "role": "customer"
  }
}
```

## `user.created` event

After the user is successfully stored, the API publishes a `user.created` event to the existing topics `finance-api` and `akademia-api`. When `X-Request-Origin: barber-app` is supplied, it also publishes that event to the `barber-app` topic. Failed user creation does not publish any event.

The event data object contains:

```json
{
  "user_id": "<created-user-id>",
  "name": "<user-name>",
  "email": "<user-email>",
  "role": "customer"
}
```

`role` is `customer` for ordinary registrations and `barber` only for a request authorized by a valid admin JWT. The `barber-app` consumer uses these fields to create the corresponding customer or barber profile.

## Authentication and identity contract

`POST /v1/user/login` keeps the existing `token` HttpOnly cookie session. Its JSON response contains a safe `user` object with `id`, `name`, `email`, `role`, `avatar`, `created_at`, and `updated_at`. Password hashes, auth keys, and payment-provider identifiers are never included in this response. The same safe shape is used by authenticated user profile responses. `GET /v1/user/me` returns the current authenticated user's profile and accepts the session cookie or an `Authorization: Bearer <token>` header.

The signed RS256 token has `sub` and `id` claims containing the same stable user UUID, plus `email`, `role`, `iss`, `iat`, and `exp`. `id` remains available for compatibility; consumers should treat `sub` as the standard subject and may use the matching `id` claim. Token roles are always one of `admin`, `barber`, or `customer`. Existing records with the historical `user` role are exposed as `customer`.

Frontend and Barber API consumers should use `user.id` as the stable identity and `user.role` for role-aware behavior. Barber API server code should validate the token signature and expiry with the configured RSA public key before trusting `sub`/`id` or `role`; never authorize from a browser-supplied role or user ID. Public signup creates a customer. A barber account can be requested only with an admin-authenticated request; public signup cannot assign the `admin` role. Admin accounts must be provisioned through a trusted administrative process.

Authorization for user records is scoped as follows:

- `GET /v1/user/me` requires a valid session or bearer token.
- `GET /v1/user` and `DELETE /v1/user/:id` require an admin token.
- `GET /v1/user/:id`, `GET /v1/user/account/:email`, and `PUT /v1/user/:id` require the account owner or an admin. Profile updates accept only `name`, `email`, and `avatar`; they cannot change roles, credentials, auth keys, or payment identifiers.
- `GET /v1/user/external/:id` requires an admin token and returns only the safe user profile.
