# Playlist publication transaction mismatch

## Bug

A success notification and audit entry are published before the playlist transaction commits. The rollback error replaces the commit error, and an attempt-specific event key lets a later delivery of the same success event produce another notification.

## Trigger

Publish `playlist-a`. The repository deterministically fails its first commit and succeeds on the second attempt, then deliver the same successful publication again with the next attempt number.

## Error

The playlist eventually becomes visible, but three success notifications are stored, the first two are premature, and the failed audit retains only the rollback error between two success entries.
