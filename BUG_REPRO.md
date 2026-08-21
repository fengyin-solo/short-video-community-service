# Feed recipient snapshot corruption

## Bug

Refreshing followers while a video notification is being aggregated changes the in-flight recipients and the cached historical recipients. Mutating a cached result can also change the current follower list.

## Trigger

Start aggregating a notification for `video-a`, wait until the background aggregation is reading its recipients, replace the first follower, then release the aggregation. Read and mutate the cached recipients afterward.

## Error

The race detector reports concurrent access to the recipient slice. The test also reports `in-flight notification was reassigned: [follower-b follower-c]`, `cached recipients changed after refresh: [follower-b follower-c]`, and `cached result mutated current followers: [intruder follower-c]`.
