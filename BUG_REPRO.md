# Safety scan cancellation leak

## Bug

A temporary safety-scan failure schedules a retry that survives request cancellation. After the request returns, the retry performs another scan and sends its late result to a closed channel.

## Trigger

Start a video safety scan, wait for the first temporary failure to schedule a retry, cancel the request and wait for the cancellation response, then release the queued retry.

## Error

The scanner call count grows from one to two after the request returns. The worker also reports `late scan result hit a closed channel: send on closed channel`.
