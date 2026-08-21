# Subtitle batch resource and transaction failure

## Bug

Subtitle validation defers source cleanup until the whole loop exits, exhausting the source limit. A malformed import then has its parse error replaced by a successful deferred commit, exposing partial subtitles and recording a successful audit.

## Trigger

Validate three normal subtitle files with a source limit of two. Then import a batch containing one normal file followed by a malformed file.

## Error

The third validation reports `too many subtitle sources open`. The malformed import returns nil, makes the first cue visible, and records `succeeded`.
