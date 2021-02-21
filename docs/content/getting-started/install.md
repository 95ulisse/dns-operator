# Installing dns-operator

`dns-operator` runs in your cluster a regular deployment.
It utilizes [CustomResourceDefinitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
to provide the new `DNSProvider` and `DNSResource` that you will later use.

It is deployed with regular Kubernetes YAML manifests.

## Installing the manifests

To install `dns-operator` in your cluster, use the all-in-one manifest.
This will install the new CRDs and correctly setup RBAC roles.

```raw
$ kubectl apply -f https://raw.githubusercontent.com/95ulisse/dns-operator/master/config/release/latest/all-in-one.yaml
```

## Verifying the installation

If the deployment went smoothly, you should see a single `dns-operator` pod marked as `Ready`.

```raw
$ kubectl get deployments/dns-operator -n dns-operator

NAME           READY   UP-TO-DATE   AVAILABLE   AGE
dns-operator   1/1     1            1           1m
```

## Configuring your first provider

Before you can start managing your DNS records, you need to configure at least one `DNSProvider`.
Go to the [Quick Start](/getting-started/quick-start) section to learn how to configure one.