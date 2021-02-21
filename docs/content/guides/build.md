# Build from source

!!! note
    `dns-operator` is built using [Kubebuilder](https://kubebuilder.io).

Clone the repository and use `make` to build the source.

```sh
$ git clone https://github.com/95ulisse/dns-operator.git
$ cd dns-operator
$ make
```

You can then run the operator either out-of-cluster or in-cluster:

- **out-of-cluster**, by directly running the compiled binary (it will connect to the cluster configured for your `kubectl`).
  This is the recommended method, as it is way simpler and faster for debugging.

    ```sh
    $ make run ENABLE_WEBHOOKS=false
    ```

- **in-cluster**, by building a Docker image and deploying the operator as a workload on your cluster.

    ```sh
    $ make docker-build IMG=custom/image-name:tag
    $ make docker-push IMG=custom/image-name:tag
    $ make deploy # This automatically installs CRDs and configures RBAC for dns-operator
    ```

## Building the docs

!!! note
    This documentation is written using [MkDocs](https://www.mkdocs.org) and [MkDocs Material](https://squidfunk.github.io/mkdocs-material).

Go to the `docs` directory and spin up a local server to test documentation.

```sh
$ cd docs
$ make docs-serve # Docs will be available on http://localhost:8000/
```