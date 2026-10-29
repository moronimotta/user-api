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
  "user_id": "<created-user-id>"
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
