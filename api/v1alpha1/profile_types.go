/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// HelmChartKind defines properties about the underlying helm chart for an artifact.
const HelmChartKind = "HelmChart"

// KustomizeKind defines a kind containing kustomize yaml files for an artifact.
const KustomizeKind = "Kustomize"

// ProfileKind defines the kind of a profile artifact
const ProfileKind = "Profile"

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// NOTE: Run "make" to regenerate code after modifying this file

// ProfileDefinitionSpec defines the desired state of ProfileDefinition
type ProfileDefinitionSpec struct {
	// Description is some text to allow a user to identify what this profile installs.
	Description string `json:"description,omitempty"`
	// Artifacts is a list of Profile artifacts
	Artifacts []Artifact `json:"artifacts,omitempty"`
}

// Artifact defines a bundled resource of the components for this profile.
type Artifact struct {
	// Name is the name of the Artifact
	Name string `json:"name,omitempty"`
	// Path is the local path to the Artifact in the Profile repo.
	// This is an optional value. If defined, it takes precedence over Chart.
	// +optional
	Path string `json:"path,omitempty"`
	// Kind is the kind of artifact: HelmChart or Kustomize
	// +kubebuilder:validation:Enum=HelmChart;Kustomize
	Kind string `json:"kind,omitempty"`
	// Chart defines properties to access a remote chart.
	// This is an optional value. It is ignored in case Path is defined.
	// +optional
	Chart *Chart `json:"chart,omitempty"`
	// Profiles defines properties to access a remote profile.
	// +optional
	Profile *Profile `json:"profile,omitempty"`
}

// Validate will validate Artifacts properties.
func (in Artifact) Validate() error {
	if in.Chart != nil && in.Path != "" {
		return apis.ErrMultipleOneOf("chart", "path")
	}
	if in.Chart != nil && in.Profile != nil {
		return apis.ErrMultipleOneOf("chart", "profile")
	}
	if in.Profile != nil && in.Path != "" {
		return apis.ErrMultipleOneOf("profile", "path")
	}
	return nil
}

// Chart defines properties to access remote helm charts.
type Chart struct {
	// URL is the URL of the Helm repository containing a Helm chart and possible values
	URL string `json:"url,omitempty"`
	// Name defines the name of the chart at the remote repository
	Name string `json:"name,omitempty"`
	// Version defines the version of the chart at the remote repository
	Version string `json:"version,omitempty"`
}

// Profile defines properties for accessing a profile
type Profile struct {
	// URL is the URL of the profile
	URL string `json:"url,omitempty"`
	// Branch is the branch in the git repository the profile lives in
	Branch string `json:"branch,omitempty"`
	// Path is the location in the git repo containing the profile definition. Only used in combination with Branch
	// +optional
	Path string `json:"path,omitempty"`
	// Version is the git tag containing the profile definition
	// +optional
	Version string `json:"version,omitempty"`
}

// ProfileDefinitionStatus defines the observed state of ProfileDefinition
// This is not used
type ProfileDefinitionStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProfileDefinition is the Schema for the profiles API
type ProfileDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileDefinitionSpec   `json:"spec,omitempty"`
	Status ProfileDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileDefinitionList contains a list of ProfileDefinition
type ProfileDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileDefinition{}, &ProfileDefinitionList{})
}
