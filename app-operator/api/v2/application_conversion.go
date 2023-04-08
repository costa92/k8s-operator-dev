package v2

import (
	dv1 "github.com/costa92/app-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Application to the Hub version (v1).
func (src *Application) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*dv1.Application)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Deployment = src.Spec.Workflow
	dst.Spec.Service = src.Spec.Service

	dst.Status.Workflow = src.Status.Workflow
	dst.Status.Network = src.Status.Network

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *Application) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dv1.Application)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Workflow = src.Spec.Deployment
	dst.Spec.Service = src.Spec.Service

	dst.Status.Workflow = src.Status.Workflow
	dst.Status.Network = src.Status.Network

	return nil
}
