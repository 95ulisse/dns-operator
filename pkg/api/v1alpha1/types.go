package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ipv4String is a string containing an IPv4 address.
// +kubebuilder:validation:Format=ipv4
type Ipv4String string

// Ipv6String is a string containing an IPv6 address.
// +kubebuilder:validation:Format=ipv6
type Ipv6String string

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
	Namespace *string `json:"namespace,omitempty"`
}

// SecretReference is a reference to a specific secret.
type SecretReference struct {
	// The name of the Secret resource being referred to.
	ObjectReference `json:",inline"`

	// The key of the entry in the Secret resource's `data` field to be used.
	Key string `json:"key,omitempty"`
}

// ConditionType enumerates the possible values of the field `Type` of a condition.
type ConditionType string

const (
	// ReadyCondition represents the `Ready` condition.
	ReadyCondition ConditionType = "Ready"
)

// ConditionStatus represents the possible values of a condition: True, False or Unknown.
type ConditionStatus string

// ConditionStatus represents the possible values of a condition: True, False or Unknown.
const (
	TrueStatus    ConditionStatus = "True"
	FalseStatus   ConditionStatus = "False"
	UnknownStatus ConditionStatus = "Unknown"
)

// Condition represents the state of a resource at a certain point in time.
// Examples of conditions are `Ready` or `Succeeded`.
type Condition struct {
	Type               ConditionType   `json:"type"`
	Status             ConditionStatus `json:"status"`
	Reason             string          `json:"reason,omitempty"`
	Message            string          `json:"message,omitempty"`
	LastTransitionTime metav1.Time     `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     metav1.Time     `json:"lastUpdateTime,omitempty"`
}

// StatusWithConditions marks a status subresource which exposes a list of conditions.
type StatusWithConditions struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

// SetCondition sets the value of a condition on a status subresource.
func (status *StatusWithConditions) SetCondition(condition *Condition) {

	if existing := status.GetCondition(condition.Type); existing >= 0 {

		// If the value of the condition changed, update the transition time
		found := &status.Conditions[existing]
		found.LastUpdateTime = metav1.Now()
		if found.Status != condition.Status || found.Reason != condition.Reason || found.Message != condition.Message {
			found.LastTransitionTime = found.LastUpdateTime
		}
		found.Status = condition.Status
		found.Reason = condition.Reason
		found.Message = condition.Message

	} else {

		// Append the condition to the list
		condition.LastUpdateTime = metav1.Now()
		condition.LastTransitionTime = condition.LastUpdateTime
		status.Conditions = append(status.Conditions, *condition)

	}

}

// GetCondition returns the index of the condition of type `conditionType` included in the given status subresource.
// Returns -1 if no condition is found.
func (status *StatusWithConditions) GetCondition(conditionType ConditionType) int {

	// Look for a condition of the same type
	for i, c := range status.Conditions {
		if c.Type == conditionType {
			return i
		}
	}
	return -1

}
