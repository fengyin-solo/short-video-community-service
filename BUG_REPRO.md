# Video delivery batch lifecycle failure

## Bug

A delivery batch never finishes after a duplicate region is skipped, even when the other region has completed. A later timed-out worker can then send to a result channel that the service already closed.

## Trigger

Deliver one duplicate region and one immediately available region in the same batch. Next, start a slow-region batch, let it time out, and release the slow worker only after the timeout returns.

## Error

The first batch reports `completed delivery batch timed out: video delivery timed out`. The late worker reports `late delivery sent to a closed result channel: send on closed channel`.
