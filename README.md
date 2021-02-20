# dns-operator

`dns-operator` is an operator for [Kubernetes](https://kubernetes.io) which allows you
to deploy DNS records as resources in you cluster.

`dns-operator` introduces a new [`DNSRecord`](https://95ulisse.github.io/dns-operator/reference/dnsrecord) kind,
which represents a single DNS record on the provider of your choice (e.g., Cloudflare or you own personal server)
that you can deploy with `kubectl`, `helm` or any other tool you use for your deployments.

For a quick overview of the core concepts of the operator, go to the [Quick Start](https://95ulisse.github.io/dns-operator/getting-started/quick-start) page,
or, if you prefer a complete step-by-step guide, head over to one of the [User Guides](https://95ulisse.github.io/dns-operator/guides/expose-an-application).

> Full documentation at: [https://95ulisse.github.io/dns-operator](https://95ulisse.github.io/dns-operator).

## License

This project is licensed under the [MIT license](LICENSE).