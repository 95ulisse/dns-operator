module github.com/95ulisse/dns-operator

go 1.13

require (
	github.com/cloudflare/cloudflare-go v0.13.8
	github.com/go-logr/logr v0.1.0
	github.com/miekg/dns v1.1.35
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)
