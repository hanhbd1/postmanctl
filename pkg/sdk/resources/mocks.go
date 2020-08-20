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

// MockListResponse represents the top-level mocks response from the
// Postman API.
type MockListResponse struct {
	Mocks MockListItems `json:"mocks,omitempty"`
}

// MockListItems is a slice of MockListItem.
type MockListItems []MockListItem

// Format returns column headers and values for the resource.
func (r MockListItems) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"UID", "Name"}, s
}

// MockListItem represents a mock in a list of all mocks.
type MockListItem struct {
	ID          string     `json:"id,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	UID         string     `json:"uid,omitempty"`
	Collection  string     `json:"collection,omitempty"`
	MockURL     string     `json:"mockUrl,omitempty"`
	Name        string     `json:"name,omitempty"`
	Config      MockConfig `json:"config,omitempty"`
	Environment string     `json:"environment,omitempty"`
}

// MockResponse is the top-level mock response from the
// Postman API.
type MockResponse struct {
	Mock Mock `json:"mock,omitempty"`
}

// Mock represents a representation of a mock server from the Postman API.
type Mock struct {
	ID          string     `json:"id,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	UID         string     `json:"uid,omitempty"`
	Collection  string     `json:"collection,omitempty"`
	MockURL     string     `json:"mockUrl,omitempty"`
	Name        string     `json:"name,omitempty"`
	Config      MockConfig `json:"config,omitempty"`
	Environment string     `json:"environment,omitempty"`
}

// MockConfig represents the configuration of a mock server.
type MockConfig struct {
	Headers          []interface{} `json:"headers,omitempty"`
	MatchBody        bool          `json:"matchBody,omitempty"`
	MatchQueryParams bool          `json:"matchQueryParams,omitempty"`
	MatchWildcards   bool          `json:"matchWildcards,omitempty"`
}

// Format returns column headers and values for the resource.
func (r Mock) Format() ([]string, []interface{}) {
	s := make([]interface{}, 1)
	s[0] = r

	return []string{"UID", "Name"}, s
}

// MockSlice is a slice of Mock.
type MockSlice []*Mock

// Format returns column headers and values for the resource.
func (r MockSlice) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"UID", "Name"}, s
}
