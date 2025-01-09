/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DirectoryCleanerSpec defines the desired state of DirectoryCleaner
type DirectoryCleanerSpec struct {
	// Paths specifies the directories/files to clean on worker nodes
	Paths []string `json:"paths"`

	// Recursive specifies whether to clean directories recursively
	Recursive bool `json:"recursive,omitempty"`

	// Force specifies whether to force clean read-only files
	Force bool `json:"force,omitempty"`

	// NodeSelector specifies which nodes should perform the cleaning
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// DirectoryCleanerStatus defines the observed state of DirectoryCleaner
type DirectoryCleanerStatus struct {
	// LastCleanedTime is the last time the cleaning operation was performed
	LastCleanedTime metav1.Time `json:"lastCleanedTime,omitempty"`

	// LastCleanedPaths contains the paths that were cleaned in the last operation
	LastCleanedPaths []string `json:"lastCleanedPaths,omitempty"`

	// Conditions represents the current state of the DirectoryCleaner
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DirectoryCleaner is the Schema for the directorycleaners API
type DirectoryCleaner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DirectoryCleanerSpec   `json:"spec,omitempty"`
	Status DirectoryCleanerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DirectoryCleanerList contains a list of DirectoryCleaner
type DirectoryCleanerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DirectoryCleaner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DirectoryCleaner{}, &DirectoryCleanerList{})
}
