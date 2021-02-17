// +build !ignore_autogenerated

/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/95ulisse/dns-operator/pkg/dnsname"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProvider) DeepCopyInto(out *DNSProvider) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProvider.
func (in *DNSProvider) DeepCopy() *DNSProvider {
	if in == nil {
		return nil
	}
	out := new(DNSProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSProvider) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProviderCloudflare) DeepCopyInto(out *DNSProviderCloudflare) {
	*out = *in
	if in.Email != nil {
		in, out := &in.Email, &out.Email
		*out = new(string)
		**out = **in
	}
	if in.APITokenSecretRef != nil {
		in, out := &in.APITokenSecretRef, &out.APITokenSecretRef
		*out = new(SecretReference)
		(*in).DeepCopyInto(*out)
	}
	if in.APIKeySecretRef != nil {
		in, out := &in.APIKeySecretRef, &out.APIKeySecretRef
		*out = new(SecretReference)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProviderCloudflare.
func (in *DNSProviderCloudflare) DeepCopy() *DNSProviderCloudflare {
	if in == nil {
		return nil
	}
	out := new(DNSProviderCloudflare)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProviderList) DeepCopyInto(out *DNSProviderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DNSProvider, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProviderList.
func (in *DNSProviderList) DeepCopy() *DNSProviderList {
	if in == nil {
		return nil
	}
	out := new(DNSProviderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSProviderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProviderRFC2136) DeepCopyInto(out *DNSProviderRFC2136) {
	*out = *in
	if in.TSIGSecretRef != nil {
		in, out := &in.TSIGSecretRef, &out.TSIGSecretRef
		*out = new(SecretReference)
		(*in).DeepCopyInto(*out)
	}
	if in.TSIGKeyName != nil {
		in, out := &in.TSIGKeyName, &out.TSIGKeyName
		*out = new(string)
		**out = **in
	}
	if in.TSIGAlgorithm != nil {
		in, out := &in.TSIGAlgorithm, &out.TSIGAlgorithm
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProviderRFC2136.
func (in *DNSProviderRFC2136) DeepCopy() *DNSProviderRFC2136 {
	if in == nil {
		return nil
	}
	out := new(DNSProviderRFC2136)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProviderSpec) DeepCopyInto(out *DNSProviderSpec) {
	*out = *in
	if in.Zones != nil {
		in, out := &in.Zones, &out.Zones
		*out = make([]dnsname.Name, len(*in))
		copy(*out, *in)
	}
	if in.Dummy != nil {
		in, out := &in.Dummy, &out.Dummy
		*out = new(bool)
		**out = **in
	}
	if in.RFC2136 != nil {
		in, out := &in.RFC2136, &out.RFC2136
		*out = new(DNSProviderRFC2136)
		(*in).DeepCopyInto(*out)
	}
	if in.Cloudflare != nil {
		in, out := &in.Cloudflare, &out.Cloudflare
		*out = new(DNSProviderCloudflare)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProviderSpec.
func (in *DNSProviderSpec) DeepCopy() *DNSProviderSpec {
	if in == nil {
		return nil
	}
	out := new(DNSProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSProviderStatus) DeepCopyInto(out *DNSProviderStatus) {
	*out = *in
	in.StatusWithConditions.DeepCopyInto(&out.StatusWithConditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSProviderStatus.
func (in *DNSProviderStatus) DeepCopy() *DNSProviderStatus {
	if in == nil {
		return nil
	}
	out := new(DNSProviderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecord) DeepCopyInto(out *DNSRecord) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecord.
func (in *DNSRecord) DeepCopy() *DNSRecord {
	if in == nil {
		return nil
	}
	out := new(DNSRecord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSRecord) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordList) DeepCopyInto(out *DNSRecordList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DNSRecord, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordList.
func (in *DNSRecordList) DeepCopy() *DNSRecordList {
	if in == nil {
		return nil
	}
	out := new(DNSRecordList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSRecordList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordSetData) DeepCopyInto(out *DNSRecordSetData) {
	*out = *in
	if in.A != nil {
		in, out := &in.A, &out.A
		*out = make([]Ipv4String, len(*in))
		copy(*out, *in)
	}
	if in.AAAA != nil {
		in, out := &in.AAAA, &out.AAAA
		*out = make([]Ipv6String, len(*in))
		copy(*out, *in)
	}
	if in.MX != nil {
		in, out := &in.MX, &out.MX
		*out = make([]MXRData, len(*in))
		copy(*out, *in)
	}
	if in.CNAME != nil {
		in, out := &in.CNAME, &out.CNAME
		*out = make([]dnsname.Name, len(*in))
		copy(*out, *in)
	}
	if in.TXT != nil {
		in, out := &in.TXT, &out.TXT
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordSetData.
func (in *DNSRecordSetData) DeepCopy() *DNSRecordSetData {
	if in == nil {
		return nil
	}
	out := new(DNSRecordSetData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordSpec) DeepCopyInto(out *DNSRecordSpec) {
	*out = *in
	in.ProviderRef.DeepCopyInto(&out.ProviderRef)
	out.Name = in.Name
	in.RRSet.DeepCopyInto(&out.RRSet)
	if in.TTLSeconds != nil {
		in, out := &in.TTLSeconds, &out.TTLSeconds
		*out = new(uint32)
		**out = **in
	}
	if in.DeletionPolicy != nil {
		in, out := &in.DeletionPolicy, &out.DeletionPolicy
		*out = new(DeletionPolicy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordSpec.
func (in *DNSRecordSpec) DeepCopy() *DNSRecordSpec {
	if in == nil {
		return nil
	}
	out := new(DNSRecordSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordStatus) DeepCopyInto(out *DNSRecordStatus) {
	*out = *in
	in.StatusWithConditions.DeepCopyInto(&out.StatusWithConditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordStatus.
func (in *DNSRecordStatus) DeepCopy() *DNSRecordStatus {
	if in == nil {
		return nil
	}
	out := new(DNSRecordStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MXRData) DeepCopyInto(out *MXRData) {
	*out = *in
	out.Host = in.Host
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MXRData.
func (in *MXRData) DeepCopy() *MXRData {
	if in == nil {
		return nil
	}
	out := new(MXRData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretReference) DeepCopyInto(out *SecretReference) {
	*out = *in
	in.ObjectReference.DeepCopyInto(&out.ObjectReference)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretReference.
func (in *SecretReference) DeepCopy() *SecretReference {
	if in == nil {
		return nil
	}
	out := new(SecretReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatusWithConditions) DeepCopyInto(out *StatusWithConditions) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusWithConditions.
func (in *StatusWithConditions) DeepCopy() *StatusWithConditions {
	if in == nil {
		return nil
	}
	out := new(StatusWithConditions)
	in.DeepCopyInto(out)
	return out
}
