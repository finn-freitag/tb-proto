# tb-proto

Protocol between the [Tossbox](https://github.com/finn-freitag/tossbox) dashboard and the [tb-agent](https://github.com/finn-freitag/tb-agent) node agent.

- `proto/tossbox/agent/v1/agent.proto`: the gRPC contract (enrollment + bidirectional control stream)
- `gen/agent/v1`: generated Go code (committed, so consumers don't need `protoc`)
- `enroll`: enrollment token format and dashboard CA pinning helpers

## Model

The dashboard owns the **desired** state of every device, and the agent owns the **actual** state:

1. The agent opens a long-lived `NodeControl.Connect` stream (outbound from the node, mTLS) and sends `Hello` with its capabilities and current devices.
2. The dashboard answers with `Welcome`, which contains the complete desired device list for this node. The agent destroys devices it manages that aren't in the list.
3. Afterwards the dashboard sends incremental `ApplyDevice` / `DeleteDevice` messages, and the agent reports `DeviceStatus` on every change plus periodic `Heartbeat`s.

## Enrollment

Token format: `tb1.<sha256 of dashboard CA cert DER>.<secret>`. The agent pins the CA by fingerprint, sends a CSR to `Enrollment.Enroll` and gets a client certificate back.

## Regenerate

```bash
nix-shell --run 'make generate'
```
