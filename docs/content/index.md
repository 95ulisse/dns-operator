# Welcome

`dns-operator` is an operator for [Kubernetes](https://kubernetes.io) which allows you
to deploy DNS records as resources in you cluster.

`dns-operator` introduces a new `DNSRecord` kind, which represents a single DNS record
on the provider of your choice (e.g., Cloudflare or you own personal server) that you can deploy
with `kubectl`, `helm` or any other tool you use for your deployments.

For a quick overview of the core concepts of the operator, go to the [Quick Start](/getting-started/quick-start) page,
or, if you prefer a complete step-by-step guide, head over to one of the [User Guides](/guides/expose-an-application).

!!! important
    **`dns-operator` is not a substitute for a complete DNS management solution!**

    This operator can help manage some of the most common record types as Kubernetes resources,
    leveraging its ecosystem of tools to ease deployments, but it is in no way a complete solution
    for all your DNS needs!