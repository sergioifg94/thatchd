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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestWorkerSpec defines the desired state of TestWorker
type TestWorkerSpec struct {
	Strategy Strategy `json:"strategy"`
}

// TestWorkerStatus defines the observed state of TestWorker
type TestWorkerStatus struct {
	DispatchedAt   *string `json:"dispatchedAt,omitempty"`
	StartedAt      *string `json:"startedAt,omitempty"`
	FinishedAt     *string `json:"finishedAt,omitempty"`
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TestWorker is the Schema for the testworkers API
type TestWorker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestWorkerSpec   `json:"spec,omitempty"`
	Status TestWorkerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestWorkerList contains a list of TestWorker
type TestWorkerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestWorker `json:"items"`
}

var _ StrategyBacked = &TestWorker{}

func (tw *TestWorker) GetStrategy() Strategy {
	return tw.Spec.Strategy
}

func init() {
	SchemeBuilder.Register(&TestWorker{}, &TestWorkerList{})
}