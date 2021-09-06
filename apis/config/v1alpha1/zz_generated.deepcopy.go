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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdvOperatorConfig) DeepCopyInto(out *AdvOperatorConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdvOperatorConfig.
func (in *AdvOperatorConfig) DeepCopy() *AdvOperatorConfig {
	if in == nil {
		return nil
	}
	out := new(AdvOperatorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdvertisementConfig) DeepCopyInto(out *AdvertisementConfig) {
	*out = *in
	out.OutgoingConfig = in.OutgoingConfig
	out.IngoingConfig = in.IngoingConfig
	if in.LabelPolicies != nil {
		in, out := &in.LabelPolicies, &out.LabelPolicies
		*out = make([]LabelPolicy, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdvertisementConfig.
func (in *AdvertisementConfig) DeepCopy() *AdvertisementConfig {
	if in == nil {
		return nil
	}
	out := new(AdvertisementConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BroadcasterConfig) DeepCopyInto(out *BroadcasterConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BroadcasterConfig.
func (in *BroadcasterConfig) DeepCopy() *BroadcasterConfig {
	if in == nil {
		return nil
	}
	out := new(BroadcasterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterConfig) DeepCopyInto(out *ClusterConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterConfig.
func (in *ClusterConfig) DeepCopy() *ClusterConfig {
	if in == nil {
		return nil
	}
	out := new(ClusterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterConfigList) DeepCopyInto(out *ClusterConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterConfigList.
func (in *ClusterConfigList) DeepCopy() *ClusterConfigList {
	if in == nil {
		return nil
	}
	out := new(ClusterConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterConfigSpec) DeepCopyInto(out *ClusterConfigSpec) {
	*out = *in
	in.AdvertisementConfig.DeepCopyInto(&out.AdvertisementConfig)
	in.DiscoveryConfig.DeepCopyInto(&out.DiscoveryConfig)
	in.LiqonetConfig.DeepCopyInto(&out.LiqonetConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterConfigSpec.
func (in *ClusterConfigSpec) DeepCopy() *ClusterConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterConfigStatus) DeepCopyInto(out *ClusterConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterConfigStatus.
func (in *ClusterConfigStatus) DeepCopy() *ClusterConfigStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscoveryConfig) DeepCopyInto(out *DiscoveryConfig) {
	*out = *in
	if in.ClusterLabels != nil {
		in, out := &in.ClusterLabels, &out.ClusterLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscoveryConfig.
func (in *DiscoveryConfig) DeepCopy() *DiscoveryConfig {
	if in == nil {
		return nil
	}
	out := new(DiscoveryConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelPolicy) DeepCopyInto(out *LabelPolicy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelPolicy.
func (in *LabelPolicy) DeepCopy() *LabelPolicy {
	if in == nil {
		return nil
	}
	out := new(LabelPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LiqonetConfig) DeepCopyInto(out *LiqonetConfig) {
	*out = *in
	if in.ReservedSubnets != nil {
		in, out := &in.ReservedSubnets, &out.ReservedSubnets
		*out = make([]CIDR, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalPools != nil {
		in, out := &in.AdditionalPools, &out.AdditionalPools
		*out = make([]CIDR, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LiqonetConfig.
func (in *LiqonetConfig) DeepCopy() *LiqonetConfig {
	if in == nil {
		return nil
	}
	out := new(LiqonetConfig)
	in.DeepCopyInto(out)
	return out
}
