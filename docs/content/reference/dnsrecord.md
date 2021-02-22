# DNSRecord Resource

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSRecord
metadata:
  name: my-record
  namespace: dns-operator
spec:

  # Full name of the DNS record.
  name: foo.example.com

  # Reference to the provider managing this record.
  providerRef:
    name: my-provider
    namespace: dns-operation # Optional, defaults to the same namespace of the DNSRecord itself

  # Specifies how to treat deletion of this DNSRecord.
  # Valid values are:
  # - "Delete" (default): actually delete the corresponding DNS record managed by this resource.
  # - "Retain": keep the published DNS record even after this resource is deleted.
  deletionPolicy: Delete

  # TTL in seconds of the DNS record. Defaults to 1h.
  ttlSeconds: 3600

  # Actual contents of the record.
  # **Only one of the nested fields can be present.**
  # Each entry in any of the nested arrays represents a *single* DNS record.
  rrset:
    a:
      - 1.1.1.1
    aaaa:
      - 1.1.1.1
    mx:
      - preference: 10
        host: mail.example.com
    cname:
      - mail.example.com
    txt:
      - Contents of the TXT record
```

!!! important
    **A single `DNSRecord` resource describes a whole RRset**, i.e., all the records of the same name and the same type in a zone.
    
    This means that if you register an `A` record for `foo.example.com` with `dns-operator`, then `dns-operator` expects to manage
    *all* the `A` records for `foo.example.com`.