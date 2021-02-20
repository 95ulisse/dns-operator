# Quick Start

!!! note
    This section assumes that you have already installed the operator as described in [Install](/getting-started/install).

## Configuring your first provider

Let's start by creating a new namespace to hold our test resources.

```sh
kubectl create namespace dns-operator-test
```

The first step is to configure a new `DNSProvider`.
As an example, we will configure a `DNSProvider` that uses [Cloudflare](https://www.cloudflare.com) as the backing provider.

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSProvider
metadata:
  name: cf-provider # Name of the provider
  namespace: dns-operator-test
spec:
  zones: # Zones that this provider will be responsible for
    - example.com
  cloudflare:
    apiTokenSecretRef:
      name: cf-provider-api-token # Secret containing the Cloudflare API Token
      key: token
---
apiVersion: v1
kind: Secret
data:
  token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX # Your Cloudflare API Token
metadata:
  name: cf-provider-api-token
  namespace: dns-operator-test
```

To check whether the provider we just deployed is ready to serve registration requests, check its `Ready` condition.
When the provider is marked as `Ready`, all the checks passed, and we can start deploying the records.
In case something is wrong in the configuration of the provider, an event will explain the details of the error.

```raw
$ kubectl describe dnsprovider/cf-provider -n dns-operator-test

[...]

Name:         cf-provider
Namespace:    dns-operator-test
API Version:  dns.k8s.marcocameriero.net/v1alpha1
Kind:         DNSProvider
Spec:
  Cloudflare:
    API Token Secret Ref:
      Key:        token
      Name:       cf-provider-api-token
  Zones:
    example.com
Status:
  Conditions:
    Last Transition Time:  2021-02-20T17:13:18Z
    Last Update Time:      2021-02-20T17:13:18Z
    Message:               Ready to register DNS records
    Reason:                Ready
    Status:                True
    Type:                  Ready
Events:
  Type     Reason  Age   From                        Message
  ----     ------  ----  ----                        -------
  Normal   Ready   12s   dns.k8s.marcocameriero.net  Ready to register DNS records
```

## Deploying your records

Let's say we want to deploy a new `A` record for `foo.example.com`.
Create and deploy the following `DNSRecord` to create a new `A` record for `foo.example.com`
pointing to `1.1.1.1`.

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSRecord
metadata:
  name: foo
  namespace: dns-operator-test
spec:
  providerRef:
    name: cf-provider # Use the Cloudflare provider we defined earlier
  name: foo.example.com
  rrset:
    a:
      - 1.1.1.1
```

To check whether the DNS record registration was successful or not, we can check the `Ready` condition of the `DNSRecord` resource:

```raw
$ kubectl describe dnsrecord/foo -n dns-operator-test

[...]

Name:         foo
Namespace:    dns-operator-test
API Version:  dns.k8s.marcocameriero.net/v1alpha1
Kind:         DNSRecord
Spec:
  Name:  foo.example.com
  Provider Ref:
    Name:  cf-provider
  Rrset:
    A:
      1.1.1.1
Status:
  Conditions:
    Last Transition Time:  2021-02-20T18:19:34Z
    Last Update Time:      2021-02-20T18:19:34Z
    Message:               DNS record registered
    Reason:                Ready
    Status:                True
    Type:                  Ready
Events:
  Type    Reason      Age   From                        Message
  ----    ------      ----  ----                        -------
  Normal  Registered  22s   dns.k8s.marcocameriero.net  DNS record correclty registered
```

We can also check that the records are actually visible.

```raw
$ dig +short foo.example.com
1.1.1.1
```

When you update the `DNSRecord` resource, the actual DNS records registered on your provider will be automatically kept in sync.
For example, change the spec of the `DNSRecord` we just created to:

```raw
$ kubectl edit dnsprovider/cf-provider -n dns-operator-test

[...]
spec:
  rrset:
    a:
      - 1.1.1.1
      - 2.2.2.2
```

And voilà.

```raw
$ dig +short foo.example.com
1.1.1.1
2.2.2.2
```

!!! important
    **A single `DNSRecord` resource describes a whole RRset**, i.e., all the records of the same name and the same type in a zone.
    
    This means that if you register an `A` record for `foo.example.com` with `dns-operator`, then `dns-operator` expects to manage
    *all* the `A` records for `foo.example.com`.

## Cleanup test resources

If you followed along this tutorial, remember to clean up the test resources.

```sh
kubectl delete namespace dns-operator-test
```