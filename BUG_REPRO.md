# Cover render request isolation failure

## Bug

An asynchronous cover render keeps a pointer to a pooled request context after the context has been returned. A later request can reuse and overwrite that object, changing both the first render target and its stored audit owner.

## Trigger

Start request A and hold its renderer after startup. Submit request B while A is blocked, let B finish, then release A and inspect A's result and audit entry.

## Error

The first result contains `request-b`, `owner-b`, and `video-b`; the audit stored under `request-a` also contains `owner-b` and `video-b`.
