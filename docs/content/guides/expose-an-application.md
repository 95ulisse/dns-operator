# Exposing an application

Let's say you have an application running on your cluster, and you want to expose it on the internet
with a public domain name.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: my-app:1.0
        ports:
        - containerPort: 8080
```

Create a `LoadBalancer` service to expose the app:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: my-app
```

!!! note
    We used here a `LoadBalancer` service just for simplicity. There are many other ways to get traffic inside your cluster.
    For more info about services,check the [official documentation](https://kubernetes.io/docs/concepts/services-networking/service/).

When the load balancer is ready, use `kubectl` to retrive the external IP.

```raw
$ kubectl get service my-app

NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
my-app   LoadBalancer   172.16.18.129   1.1.1.1       80:30773/TCP   1m
```

Use the external IP to setup a `DNSRecord` pointing to the load balancer.

```yaml
apiVersion: dns.k8s.marcocameriero.net/v1alpha1
kind: DNSRecord
metadata:
  name: my-app
spec:
  providerRef:
    name: my-provider
  name: my-app.example.com
  rrset:
    a:
      - 1.1.1.1
```

The app will now be reachable on:

```raw
http://my-app.example.com
```