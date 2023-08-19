/*
Copyright 2023 Costalong.

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

type WorkSpaceOperation string

const (
	WorkSpaceStart WorkSpaceOperation = "Start"
	WorkSpaceStop                     = "Stop"
)

type WorkSpacePhase string

const (
	WorkspacePhaseRunning WorkSpacePhase = "Running"
	WorkspacePhaseStopped WorkSpacePhase = "Stopped"
)

// WorkSpaceSpec defines the desired state of WorkSpace
type WorkSpaceSpec struct {
	// 表示该工作空间使用的cpu、内存和存储的规格
	Cpu     string `json:"cpu,omitempty"`
	Memory  string `json:"memory,omitempty"`
	Storage string `json:"storage,omitempty"`

	// 是一个用于描述硬件资源的字段，用于在使用kubectl查询时显示信息
	Hardware string `json:"hardware,omitempty"`
	// pod使用的镜像
	Image string `json:"image,omitempty"`
	// pod中code-server监听的端口
	Port int32 `json:"port,omitempty"`
	// 存储卷的挂载位置
	MountPath string `json:"mountPath"`
	// 要进行的操作，用于启动或者停止工作空间
	Operation WorkSpaceOperation `json:"operation,omitempty"`
}

// WorkSpaceStatus defines the observed state of WorkSpace
type WorkSpaceStatus struct {
	Phase WorkSpacePhase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`

// WorkSpace is the Schema for the workspaces API
type WorkSpace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkSpaceSpec   `json:"spec,omitempty"`
	Status WorkSpaceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WorkSpaceList contains a list of WorkSpace
type WorkSpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkSpace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkSpace{}, &WorkSpaceList{})
}
