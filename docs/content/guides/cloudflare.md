# Configuring the Cloudflare provider

To use Cloudflare, you can configure one of two different authentication credentials:

- **API Tokens**, which allow finer control and application-specific permessions to certain zones and resources.
- **API Keys**, which are globally-scoped keys that carry the same permissions of you whole account.

**API Tokens** are generally more recommendable for their higher livel of security with respect to API Keys.

## Using an API Token

Create a new API Token from you Cloudflare dashboard at **My Profile > API Tokens > Create Token**.
The minimum set of permissions required for `dns-operator` to correctly funcion are:

- Permissions:
    - `Zone` > `Zone` > `Read`
    - `Zone` > `DNS` > `Edit`

- Zone Resources:
    - `Include` > `Specific zone` > `<Name of the zone>` (or alternatively, grant access to all zones by selecting `All zones`)

Once you have your new API Token, store it in a Kubernetes secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cf-provider-api-token
type: Opaque
data:
  token: XXXXXXXXXXXXXXXXX # Base64 encoded
```

And then reference it in your `DNSProvider`:

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSProvider
metadata:
  name: cf-provider
spec:
  zones:
    # These zones must be included in the list of authorized zones for the token
    - '[...]'
  cloudflare:
    apiTokenSecretRef:
      name: cf-provider-api-token
      key: token
```

## Using an API Key

View you account key at **My Profile > API Tokens > Global API Key > View**.
Copy the key and store in in a Kubernetes secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cf-provider-api-key
type: Opaque
data:
  key: XXXXXXXXXXXXXXXXX # Base64 encoded
```

And then reference it in your `DNSProvider`:

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSProvider
metadata:
  name: cf-provider
spec:
  zones:
    - '[...]'
  cloudflare:
    email: your-cloudflare-email@example.com
    apiKeySecretRef:
      name: cf-provider-api-key
      key: key
```