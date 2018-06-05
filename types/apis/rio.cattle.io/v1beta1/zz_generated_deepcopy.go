package v1beta1

import (
	reflect "reflect"

	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*BindOptions).DeepCopyInto(out.(*BindOptions))
			return nil
		}, InType: reflect.TypeOf(&BindOptions{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Condition).DeepCopyInto(out.(*Condition))
			return nil
		}, InType: reflect.TypeOf(&Condition{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ContainerConfig).DeepCopyInto(out.(*ContainerConfig))
			return nil
		}, InType: reflect.TypeOf(&ContainerConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ContainerPrivilegedConfig).DeepCopyInto(out.(*ContainerPrivilegedConfig))
			return nil
		}, InType: reflect.TypeOf(&ContainerPrivilegedConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DeviceMapping).DeepCopyInto(out.(*DeviceMapping))
			return nil
		}, InType: reflect.TypeOf(&DeviceMapping{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*DriverConfig).DeepCopyInto(out.(*DriverConfig))
			return nil
		}, InType: reflect.TypeOf(&DriverConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*HealthConfig).DeepCopyInto(out.(*HealthConfig))
			return nil
		}, InType: reflect.TypeOf(&HealthConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*InternalStack).DeepCopyInto(out.(*InternalStack))
			return nil
		}, InType: reflect.TypeOf(&InternalStack{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Mount).DeepCopyInto(out.(*Mount))
			return nil
		}, InType: reflect.TypeOf(&Mount{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*PodConfig).DeepCopyInto(out.(*PodConfig))
			return nil
		}, InType: reflect.TypeOf(&PodConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*PortBinding).DeepCopyInto(out.(*PortBinding))
			return nil
		}, InType: reflect.TypeOf(&PortBinding{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*PrivilegedConfig).DeepCopyInto(out.(*PrivilegedConfig))
			return nil
		}, InType: reflect.TypeOf(&PrivilegedConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Service).DeepCopyInto(out.(*Service))
			return nil
		}, InType: reflect.TypeOf(&Service{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ServiceList).DeepCopyInto(out.(*ServiceList))
			return nil
		}, InType: reflect.TypeOf(&ServiceList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ServiceSpec).DeepCopyInto(out.(*ServiceSpec))
			return nil
		}, InType: reflect.TypeOf(&ServiceSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*SidecarConfig).DeepCopyInto(out.(*SidecarConfig))
			return nil
		}, InType: reflect.TypeOf(&SidecarConfig{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Space).DeepCopyInto(out.(*Space))
			return nil
		}, InType: reflect.TypeOf(&Space{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*SpaceSpec).DeepCopyInto(out.(*SpaceSpec))
			return nil
		}, InType: reflect.TypeOf(&SpaceSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*SpaceStatus).DeepCopyInto(out.(*SpaceStatus))
			return nil
		}, InType: reflect.TypeOf(&SpaceStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Stack).DeepCopyInto(out.(*Stack))
			return nil
		}, InType: reflect.TypeOf(&Stack{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*StackList).DeepCopyInto(out.(*StackList))
			return nil
		}, InType: reflect.TypeOf(&StackList{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*StackScoped).DeepCopyInto(out.(*StackScoped))
			return nil
		}, InType: reflect.TypeOf(&StackScoped{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*StackSpec).DeepCopyInto(out.(*StackSpec))
			return nil
		}, InType: reflect.TypeOf(&StackSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*StackStatus).DeepCopyInto(out.(*StackStatus))
			return nil
		}, InType: reflect.TypeOf(&StackStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Tmpfs).DeepCopyInto(out.(*Tmpfs))
			return nil
		}, InType: reflect.TypeOf(&Tmpfs{})},
		conversion.GeneratedDeepCopyFunc{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*VolumeOptions).DeepCopyInto(out.(*VolumeOptions))
			return nil
		}, InType: reflect.TypeOf(&VolumeOptions{})},
	)
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BindOptions) DeepCopyInto(out *BindOptions) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BindOptions.
func (in *BindOptions) DeepCopy() *BindOptions {
	if in == nil {
		return nil
	}
	out := new(BindOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	return
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
func (in *ContainerConfig) DeepCopyInto(out *ContainerConfig) {
	*out = *in
	out.ContainerPrivilegedConfig = in.ContainerPrivilegedConfig
	if in.CapAdd != nil {
		in, out := &in.CapAdd, &out.CapAdd
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CapDrop != nil {
		in, out := &in.CapDrop, &out.CapDrop
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Entrypoint != nil {
		in, out := &in.Entrypoint, &out.Entrypoint
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Healthcheck != nil {
		in, out := &in.Healthcheck, &out.Healthcheck
		if *in == nil {
			*out = nil
		} else {
			*out = new(HealthConfig)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Tmpfs != nil {
		in, out := &in.Tmpfs, &out.Tmpfs
		*out = make([]Tmpfs, len(*in))
		copy(*out, *in)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Mount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumesFrom != nil {
		in, out := &in.VolumesFrom, &out.VolumesFrom
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Devices != nil {
		in, out := &in.Devices, &out.Devices
		*out = make([]DeviceMapping, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerConfig.
func (in *ContainerConfig) DeepCopy() *ContainerConfig {
	if in == nil {
		return nil
	}
	out := new(ContainerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerPrivilegedConfig) DeepCopyInto(out *ContainerPrivilegedConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerPrivilegedConfig.
func (in *ContainerPrivilegedConfig) DeepCopy() *ContainerPrivilegedConfig {
	if in == nil {
		return nil
	}
	out := new(ContainerPrivilegedConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceMapping) DeepCopyInto(out *DeviceMapping) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceMapping.
func (in *DeviceMapping) DeepCopy() *DeviceMapping {
	if in == nil {
		return nil
	}
	out := new(DeviceMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DriverConfig) DeepCopyInto(out *DriverConfig) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DriverConfig.
func (in *DriverConfig) DeepCopy() *DriverConfig {
	if in == nil {
		return nil
	}
	out := new(DriverConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthConfig) DeepCopyInto(out *HealthConfig) {
	*out = *in
	if in.Test != nil {
		in, out := &in.Test, &out.Test
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthConfig.
func (in *HealthConfig) DeepCopy() *HealthConfig {
	if in == nil {
		return nil
	}
	out := new(HealthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalStack) DeepCopyInto(out *InternalStack) {
	*out = *in
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = make(map[string]Service, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalStack.
func (in *InternalStack) DeepCopy() *InternalStack {
	if in == nil {
		return nil
	}
	out := new(InternalStack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mount) DeepCopyInto(out *Mount) {
	*out = *in
	if in.BindOptions != nil {
		in, out := &in.BindOptions, &out.BindOptions
		if *in == nil {
			*out = nil
		} else {
			*out = new(BindOptions)
			**out = **in
		}
	}
	if in.VolumeOptions != nil {
		in, out := &in.VolumeOptions, &out.VolumeOptions
		if *in == nil {
			*out = nil
		} else {
			*out = new(VolumeOptions)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mount.
func (in *Mount) DeepCopy() *Mount {
	if in == nil {
		return nil
	}
	out := new(Mount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodConfig) DeepCopyInto(out *PodConfig) {
	*out = *in
	if in.StopGracePeriodSeconds != nil {
		in, out := &in.StopGracePeriodSeconds, &out.StopGracePeriodSeconds
		if *in == nil {
			*out = nil
		} else {
			*out = new(int)
			**out = **in
		}
	}
	if in.ExtraHosts != nil {
		in, out := &in.ExtraHosts, &out.ExtraHosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodConfig.
func (in *PodConfig) DeepCopy() *PodConfig {
	if in == nil {
		return nil
	}
	out := new(PodConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortBinding) DeepCopyInto(out *PortBinding) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortBinding.
func (in *PortBinding) DeepCopy() *PortBinding {
	if in == nil {
		return nil
	}
	out := new(PortBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivilegedConfig) DeepCopyInto(out *PrivilegedConfig) {
	*out = *in
	if in.PortBindings != nil {
		in, out := &in.PortBindings, &out.PortBindings
		*out = make([]PortBinding, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivilegedConfig.
func (in *PrivilegedConfig) DeepCopy() *PrivilegedConfig {
	if in == nil {
		return nil
	}
	out := new(PrivilegedConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.ServiceSpec.DeepCopyInto(&out.ServiceSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Service.
func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Service) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceList) DeepCopyInto(out *ServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Service, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceList.
func (in *ServiceList) DeepCopy() *ServiceList {
	if in == nil {
		return nil
	}
	out := new(ServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceSpec) DeepCopyInto(out *ServiceSpec) {
	*out = *in
	out.StackScoped = in.StackScoped
	in.PodConfig.DeepCopyInto(&out.PodConfig)
	in.PrivilegedConfig.DeepCopyInto(&out.PrivilegedConfig)
	if in.Sidecars != nil {
		in, out := &in.Sidecars, &out.Sidecars
		*out = make(map[string]SidecarConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	in.ContainerConfig.DeepCopyInto(&out.ContainerConfig)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSpec.
func (in *ServiceSpec) DeepCopy() *ServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SidecarConfig) DeepCopyInto(out *SidecarConfig) {
	*out = *in
	in.ContainerConfig.DeepCopyInto(&out.ContainerConfig)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SidecarConfig.
func (in *SidecarConfig) DeepCopy() *SidecarConfig {
	if in == nil {
		return nil
	}
	out := new(SidecarConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Space) DeepCopyInto(out *Space) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Space.
func (in *Space) DeepCopy() *Space {
	if in == nil {
		return nil
	}
	out := new(Space)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Space) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpaceSpec) DeepCopyInto(out *SpaceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpaceSpec.
func (in *SpaceSpec) DeepCopy() *SpaceSpec {
	if in == nil {
		return nil
	}
	out := new(SpaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpaceStatus) DeepCopyInto(out *SpaceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpaceStatus.
func (in *SpaceStatus) DeepCopy() *SpaceStatus {
	if in == nil {
		return nil
	}
	out := new(SpaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stack) DeepCopyInto(out *Stack) {
	*out = *in
	out.Namespaced = in.Namespaced
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stack.
func (in *Stack) DeepCopy() *Stack {
	if in == nil {
		return nil
	}
	out := new(Stack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Stack) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackList) DeepCopyInto(out *StackList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Stack, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackList.
func (in *StackList) DeepCopy() *StackList {
	if in == nil {
		return nil
	}
	out := new(StackList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StackList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackScoped) DeepCopyInto(out *StackScoped) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackScoped.
func (in *StackScoped) DeepCopy() *StackScoped {
	if in == nil {
		return nil
	}
	out := new(StackScoped)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackSpec) DeepCopyInto(out *StackSpec) {
	*out = *in
	if in.Templates != nil {
		in, out := &in.Templates, &out.Templates
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackSpec.
func (in *StackSpec) DeepCopy() *StackSpec {
	if in == nil {
		return nil
	}
	out := new(StackSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StackStatus) DeepCopyInto(out *StackStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StackStatus.
func (in *StackStatus) DeepCopy() *StackStatus {
	if in == nil {
		return nil
	}
	out := new(StackStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tmpfs) DeepCopyInto(out *Tmpfs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tmpfs.
func (in *Tmpfs) DeepCopy() *Tmpfs {
	if in == nil {
		return nil
	}
	out := new(Tmpfs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeOptions) DeepCopyInto(out *VolumeOptions) {
	*out = *in
	if in.DriverConfig != nil {
		in, out := &in.DriverConfig, &out.DriverConfig
		if *in == nil {
			*out = nil
		} else {
			*out = new(DriverConfig)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeOptions.
func (in *VolumeOptions) DeepCopy() *VolumeOptions {
	if in == nil {
		return nil
	}
	out := new(VolumeOptions)
	in.DeepCopyInto(out)
	return out
}
