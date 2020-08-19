/*
Copyright © 2020 Kevin Swiber <kswiber@gmail.com>

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

package resources

// WorkspaceListResponse represents the top-level workspaces response from the
// Postman API.
type WorkspaceListResponse struct {
	Workspaces WorkspaceListItems `json:"workspaces,omitempty"`
}

// WorkspaceListItems is a slice of WorkspaceListItem.
type WorkspaceListItems []WorkspaceListItem

// Format returns column headers and values for the resource.
func (r WorkspaceListItems) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"ID", "Name", "Type"}, s
}

// WorkspaceListItem represents a single item in an WorkspaceListResponse.
type WorkspaceListItem struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// WorkspaceResponse is the top-level workspace response from the
// Postman API.
type WorkspaceResponse struct {
	Workspace Workspace `json:"workspace,omitempty"`
}

// Workspace represents the single workspace response from the
// Postman API
type Workspace struct {
	ID           string                         `json:"id,omitempty"`
	Name         string                         `json:"name,omitempty"`
	Type         string                         `json:"type,omitempty"`
	Description  string                         `json:"description,omitempty"`
	Collections  []WorkspaceCollectionListItem  `json:"collections,omitempty"`
	Environments []WorkspaceEnvironmentListItem `json:"environments,omitempty"`
	Mocks        []WorkspaceMockListItem        `json:"mocks,omitempty"`
	Monitors     []WorkspaceMonitorListItem     `json:"monitors,omitempty"`
}

// Format returns column headers and values for the resource.
func (r Workspace) Format() ([]string, []interface{}) {
	s := make([]interface{}, 1)
	s[0] = r

	return []string{"ID", "Name", "Type"}, s
}

// WorkspaceSlice is a slice of Workspace.
type WorkspaceSlice []*Workspace

// Format returns column headers and values for the resource.
func (r WorkspaceSlice) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"ID", "Name", "Type"}, s
}

// WorkspaceCollectionListItem represents a single collection item in a Workspace.
type WorkspaceCollectionListItem struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	UID  string `json:"uid,omitempty"`
}

// WorkspaceEnvironmentListItem represents a single environment item in a Workspace.
type WorkspaceEnvironmentListItem struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	UID  string `json:"uid,omitempty"`
}

// WorkspaceMockListItem represents a single mock item in a Workspace.
type WorkspaceMockListItem struct {
	ID string `json:"id,omitempty"`
}

// WorkspaceMonitorListItem represents a single monitor item in a Workspace.
type WorkspaceMonitorListItem struct {
	ID string `json:"id,omitempty"`
}
