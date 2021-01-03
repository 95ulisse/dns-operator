package v1alpha1

// DeletionPolicy describes how the DNSRecord resource deletion will propagate to the underlying actual DNS record.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is Delete.
// +kubebuilder:validation:Enum=Delete;Retain
type DeletionPolicy string

const (
	// DeletePolicy propagates the deletion of the DNSRecord resource to the underlying actual DNS record.
	DeletePolicy DeletionPolicy = "Delete"

	// RetainPolicy does not delete the actual DNS record when the DNSRecord resource is deleted.
	RetainPolicy DeletionPolicy = "Retain"
)

// ObjectReference is a reference to an object in a (possibly another) namespace.
type ObjectReference struct {
	// Name of the resource being referred.
	Name string `json:"name"`

	// Name of the namespace of the resource being referred.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SecretReference is a reference to a specific secret.
type SecretReference struct {
	// The name of the Secret resource being referred to.
	ObjectReference `json:",inline"`

	// The key of the entry in the Secret resource's `data` field to be used.
	// Some instances of this field may be defaulted, in others it may be
	// required.
	// +optional
	Key string `json:"key,omitempty"`
}
