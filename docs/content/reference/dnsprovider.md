# DNSProvider Resource

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSProvider
metadata:
  name: my-provider
  namespace: dns-operator
spec:

  # DNS zones handled by this provider.
  # At least one zone must be present.
  zones:
    - example.com
  
  # Cloudflare provider configuration
  cloudflare:

    # If true, marks all records as proxied by default.
    # Defaults to true.
    proxiedByDefault: true

    # Email of your Cloudflare account. Required only if using an API Key.
    email: my-cloudflare-email@example.com

    # Reference to a secret containing the API Key to use for authentication.
    # Either this field or `apiTokenSecretRef` is required.
    apiKeySecretRef:
      name: cf-provider-api-key
      key: key
    
    # Reference to a secret containing the API Token to use for authentication.
    # Either this field or `apiKeySecretRef` is required.
    apiTokenSecretRef:
      name: cf-provider-api-token
      key: token

  # RFC2136 (aka Dynamic DNS) provider
  rfc2136:

    # The IP address or hostname of an authoritative DNS server supporting RFC2136 in the form host:port.
    nameserver: 1.1.1.1

    # The name of the secret containing the TSIG value.
    # If any of the `tsig*` fields is defined, this field is required.
    tsigSecretRef:
      name: my-provider-tsig-secret
      key: key

    # The TSIG Key name configured in the DNS.
    # If any of the `tsig*` fields is defined, this field is required.
    tsigKeyName: dns-operator

    # The TSIG Algorithm configured in the DNS supporting RFC2136.
    # Used only when ``tsigSecretSecretRef`` and ``tsigKeyName`` are defined.
    # Supported values are (case-insensitive):
    # "HMACMD5", "HMACSHA1", "HMACSHA256" or "HMACSHA512".
    tsigAlgorithm: HMACSHA512
```