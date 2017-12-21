Required environment variables

- `CF_SYSTEM_DOMAIN`
- `CF_USERNAME`
- `CF_PASSWORD`
- `CF_SPACE`
- `CF_ORG`
- `CYCLES`
- `DELAY_US`
- `DRAIN_URLS` a string with URLs separated by spaces
- `DATADOG_API_KEY`
- `DRAIN_VERSION`
- `SINK_DEPLOY`
- `NUM_APPS`

Required if `SINK_DEPLOY` is `standalone`,
- `EXTERNAL_DRAIN_HOST`
- `EXTERNAL_DRAIN_PORT`
- `EXTERNAL_COUNTER_PORT`

Usage:

<!-- TODO: add generate id script -->
```
./push.sh && ./hammer.sh && ./report.sh && ./teardown.sh
```
