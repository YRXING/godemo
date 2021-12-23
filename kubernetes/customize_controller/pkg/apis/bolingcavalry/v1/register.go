package v1

import (
	"godemo/kubebuilder/customize_controller/pkg/apis/bolingcavalry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{
	Group: bolingcavalry.GroupName,
	Version: bolingcavalry.Version,
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme = SchemeBuilder.AddToScheme
)

func Resouce(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func addKnownTypes(schem *runtime.Scheme) error {
	schem.AddKnownTypes(
		SchemeGroupVersion,
		&Student{},
		&StudentList{},
		)

	// register the type in the scheme
	metav1.AddToGroupVersion(schem, SchemeGroupVersion)
	return nil
}
