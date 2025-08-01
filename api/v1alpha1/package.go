// +k8s:deepcopy-gen=package
// +groupName=nokamoto.github.com
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var SchemeBuilder = &scheme.Builder{
	GroupVersion: schema.GroupVersion{
		Group:   "nokamoto.github.com",
		Version: "v1alpha1",
	},
}
