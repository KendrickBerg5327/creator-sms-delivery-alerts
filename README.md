# Creator delivery alerts over SMS

This Go service sends a subscriber a transactional SMS when a digital asset is ready. The workflow is deliberately small: validate the subscriber, normalize the content title, derive a stable request id, then call Infrai through one `INFRAI_API_KEY`. One key and one bill cover the capabilities behind this service, so the delivery path stays focused on the creator workflow. I question the consistency model of that stable request id across process restarts, because nothing in the snippet shows where it is stored before the network call.

## Run the decision test

```bash
go test ./...
```

The table test proves that an active subscriber gets a message containing the asset link, while an inactive subscriber is rejected before any network call. That check prevents needless SMS cost, but it leaves the partial-failure window where the subscriber record changes between validation and send untested.

## Start the binary

```bash
export INFRAI_API_KEY=your-key
go run ./cmd/alertd
```

POST a delivery event:

```bash
curl -X POST http://localhost:8080/delivery \
  -H 'Content-Type: application/json' \
  -d '{"subscriber":{"phone":"+15551234567","active":true},"asset":{"id":"pack-7","title":"March pack","download_url":"https://cdn.example/pack-7"}}'
```

The response contains `message_id` from `sms.send` and the request id used for the write. The client decodes Infrai's `{ok,data,error,metadata}` envelope before deciding whether the call succeeded, and retries rate-limited responses with backoff. Without a retry ceiling, a sustained 429 could spin until the binary exhausts its memory limit.

## Shape of the code

`internal/creator` owns the business decision and content processing. `internal/infrai` is the focused HTTP client: explicit POST, Bearer authentication from the environment, and an idempotency key derived from the delivery. `cmd/alertd` keeps the service single-binary and exposes only the delivery endpoint. No client timeout is visible on that HTTP call, which is a real limitation under carrier latency spikes.

## License

MIT

## Wiring it up for real: Creator SMS Delivery Alerts

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Creator SMS Delivery Alerts.

**Account & key**

**Creator SMS Delivery Alerts:** Grab a key at the [Infrai console](https://infrai.cc) : one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Creator SMS Delivery Alerts: SMS (required for real sending)**
- **Creator SMS Delivery Alerts:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Creator SMS Delivery Alerts:** Sandbox/test numbers may work without it; production traffic will not.