# Unknown playback profile failure

## Bug

An unknown device is allowed to request 4K playback. Enabling the AV1 fallback on the same profile then returns an internal error and leaves the codec disabled.

## Trigger

Load a playback profile without a device identifier, check 4K playback, then check the 2K fallback and enable AV1 on that profile.

## Error

The test reports `unknown device was allowed 4K playback`, `fallback codec update returned an internal error`, and `fallback AV1 codec was not enabled`.
