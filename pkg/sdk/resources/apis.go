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

import "time"

// APIListResponse represents the top-level APIs response from the Postman API.
type APIListResponse struct {
	APIs APIListItems `json:"apis,omitempty"`
}

// APIListItems is a slice of APIListItem
type APIListItems []APIListItem

// Format returns column headers and values for the resource.
func (r APIListItems) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"ID", "Name"}, s
}

// APIListItem represents a single item in an APIListResponse.
type APIListItem struct {
	CreatedBy   string    `json:"createdBy,omitempty"`
	UpdatedBy   string    `json:"updatedBy,omitempty"`
	Team        string    `json:"team,omitempty"`
	ID          string    `json:"-,omitempty"`
	Name        string    `json:"name,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// APIResponse is a single API representation in the Postman API.
type APIResponse struct {
	API API `json:"api,omitempty"`
}

// API represents a single item in an APIListResponse.
type API struct {
	CreatedBy   string    `json:"createdBy,omitempty"`
	UpdatedBy   string    `json:"updatedBy,omitempty"`
	Team        string    `json:"team,omitempty"`
	ID          string    `json:"-,omitempty"`
	Name        string    `json:"name,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// Format returns column headers and values for the resource.
func (r API) Format() ([]string, []interface{}) {
	s := make([]interface{}, 1)
	s[0] = r

	return []string{"ID", "Name"}, s
}

// APISlice is a slice of API
type APISlice []*API

// Format returns column headers and values for the resource.
func (r APISlice) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"ID", "Name"}, s
}
