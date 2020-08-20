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

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kevinswiber/postmanctl/pkg/util"
	"github.com/spf13/cast"
	"os"
	"strings"

	"github.com/kevinswiber/postmanctl/pkg/sdk/resources"
	"github.com/spf13/cobra"
)

var outputFormat OutputFormatValue
var outputFile defaultValue
var ignoreKey defaultValue
var removeNil defaultValue

type defaultValue struct {
	value string
}

func (o *defaultValue) String() string {
	return o.value
}

// Set creates the flag value.
func (o *defaultValue) Set(v string) error {
	o.value = v
	return nil
}

// Type returns the type of this value.
func (o *defaultValue) Type() string {
	return "string"
}

// OutputFormatValue is a custom Value for the output flag that validates.
type OutputFormatValue struct {
	value string
}

// String returns a string representation of this flag.
func (o *OutputFormatValue) String() string {
	return o.value
}

// Set creates the flag value.
func (o *OutputFormatValue) Set(v string) error {
	if v == "json" || strings.HasPrefix(v, "jsonpath=") ||
		strings.HasPrefix(v, "go-template-file=") {
		o.value = v
		return nil
	}

	return errors.New("output format must be json, jsonpath, or go-template-file")
}

// Type returns the type of this value.
func (o *OutputFormatValue) Type() string {
	return "string"
}

func init() {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve Postman resources.",
	}

	apiVersionsCmd := &cobra.Command{
		Use:     "api-versions",
		Aliases: []string{"api-version"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				params := append([]string{forAPI}, args...)
				return getIndividualAPIVersions(params)
			}

			return getAllResources(resources.APIVersionType, forAPI)
		},
	}

	apiVersionsCmd.Flags().StringVar(&forAPI, "for-api", "", "the associated API ID (required)")
	apiVersionsCmd.MarkFlagRequired("for-api")

	schemaCmd := &cobra.Command{
		Use: "schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := []string{forAPI, forAPIVersion}
			if len(args) == 0 {
				version, err := service.APIVersion(context.Background(), forAPI, forAPIVersion)

				if err != nil {
					return handleResponseError(err)
				}

				if len(version.Schema) > 0 {
					params = append(params, version.Schema[0])
				} else {
					fmt.Fprintln(os.Stderr, "error: no schema has been associated with this API version")
					os.Exit(1)
				}
			}
			params = append(params, args...)
			return getIndividualSchema(params)
		},
	}

	schemaCmd.Flags().StringVar(&forAPI, "for-api", "", "the associated API ID (required)")
	schemaCmd.MarkFlagRequired("for-api")

	schemaCmd.Flags().StringVar(&forAPIVersion, "for-api-version", "", "the associated API Version ID (required)")
	schemaCmd.MarkFlagRequired("for-api-version")

	apiRelationsCmd := &cobra.Command{
		Use: "api-relations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputFormat.value == "" {
				return getFormattedAPIRelations(forAPI, forAPIVersion)
			}
			return getAPIRelations(forAPI, forAPIVersion)
		},
	}

	apiRelationsCmd.Flags().StringVar(&forAPI, "for-api", "", "the associated API ID (required)")
	apiRelationsCmd.MarkFlagRequired("for-api")

	apiRelationsCmd.Flags().StringVar(&forAPIVersion, "for-api-version", "", "the associated API Version ID (required)")
	apiRelationsCmd.MarkFlagRequired("for-api-version")

	userCmd := &cobra.Command{
		Use: "user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getIndividualUser(args)
		},
	}

	apisCmd := generateGetSubcommand(resources.APIType, "apis", []string{"api"}, getIndividualAPIs)
	apisCmd.Flags().StringVar(&usingWorkspace, "workspace", "", "the associated workspace ID")

	getCmd.AddCommand(
		generateGetSubcommand(resources.CollectionType, "collections", []string{"collection", "co"}, getIndividualCollections),
		generateGetSubcommand(resources.EnvironmentType, "environments", []string{"environment", "env"}, getIndividualEnvironments),
		generateGetSubcommand(resources.MonitorType, "monitors", []string{"monitor", "mon"}, getIndividualMonitors),
		generateGetSubcommand(resources.MockType, "mocks", []string{"mock"}, getIndividualMocks),
		generateGetSubcommand(resources.WorkspaceType, "workspaces", []string{"workspace", "ws"}, getIndividualWorkspaces),
		userCmd,
		apisCmd,
		apiVersionsCmd,
		apiRelationsCmd,
		schemaCmd,
	)

	getCmd.PersistentFlags().VarP(&outputFormat, "output", "o", "output format (json, jsonpath, go-template-file)")
	getCmd.PersistentFlags().VarP(&outputFile, "file", "f", "output file")
	getCmd.PersistentFlags().VarP(&ignoreKey, "ignore-key", "i", "ignore json key in response")
	rootCmd.AddCommand(getCmd)
}

func generateGetSubcommand(t resources.ResourceType, use string, aliases []string, fn func(args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Aliases: aliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fn(args)
			}

			return getAllResources(t)
		},
	}
}

func prepareMap(resourceType resources.ResourceType, args ...string) map[string]string {
	ctx := context.Background()
	uuidmap := make(map[string]string)
	var resource interface{}
	var err error

	switch resourceType {
	case resources.CollectionType:
		resource, err = service.Collections(ctx)
	case resources.EnvironmentType:
		resource, err = service.Environments(ctx)
	case resources.MockType:
		resource, err = service.Mocks(ctx)
	case resources.MonitorType:
		resource, err = service.Monitors(ctx)
	case resources.APIType:
		resource, err = service.APIs(ctx, usingWorkspace)
	case resources.APIVersionType:
		resource, err = service.APIVersions(ctx, args[0])
	case resources.WorkspaceType:
		resource, err = service.Workspaces(ctx)
	default:
		return uuidmap
	}

	if err != nil {
		return uuidmap
	}
	data, _ := json.Marshal(resource)
	var tmp []map[string]interface{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return uuidmap
	}
	for _, t := range tmp {
		uuidmap[cast.ToString(t["name"])] = cast.ToString(t["uid"])
	}
	return uuidmap
}

func getAllResources(resourceType resources.ResourceType, args ...string) error {
	ctx := context.Background()

	var resource interface{}
	var err error

	switch resourceType {
	case resources.CollectionType:
		resource, err = service.Collections(ctx)
	case resources.EnvironmentType:
		resource, err = service.Environments(ctx)
	case resources.MockType:
		resource, err = service.Mocks(ctx)
	case resources.MonitorType:
		resource, err = service.Monitors(ctx)
	case resources.APIType:
		resource, err = service.APIs(ctx, usingWorkspace)
	case resources.APIVersionType:
		resource, err = service.APIVersions(ctx, args[0])
	case resources.WorkspaceType:
		resource, err = service.Workspaces(ctx)
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType.String())
	}

	if err != nil {
		return handleResponseError(err)
	}

	printGetOutput(resource)

	return nil
}

func getIndividualCollections(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.CollectionType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.Collection
		var err error
		if !ok {
			resource, err = service.Collection(context.Background(), name)
		} else {
			resource, err = service.Collection(context.Background(), id)
		}
		if err != nil {
			return handleResponseError(err)
		}
		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualEnvironments(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.EnvironmentType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.Environment
		var err error
		if !ok {
			resource, err = service.Environment(context.Background(), name)
		} else {
			resource, err = service.Environment(context.Background(), id)
		}

		if err != nil {
			return handleResponseError(err)
		}

		if err != nil {
			return handleResponseError(err)
		}
		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualMocks(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.MockType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.Mock
		var err error
		if !ok {
			resource, err = service.Mock(context.Background(), name)
		} else {
			resource, err = service.Mock(context.Background(), id)
		}

		if err != nil {
			return handleResponseError(err)
		}

		if err != nil {
			return handleResponseError(err)
		}
		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualMonitors(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.MonitorType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.Monitor
		var err error
		if !ok {
			resource, err = service.Monitor(context.Background(), name)
		} else {
			resource, err = service.Monitor(context.Background(), id)
		}

		if err != nil {
			return handleResponseError(err)
		}

		if err != nil {
			return handleResponseError(err)
		}
		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualAPIs(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.APIType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.API
		var err error
		if !ok {
			resource, err = service.API(context.Background(), name)
		} else {
			resource, err = service.API(context.Background(), id)
		}

		if err != nil {
			return handleResponseError(err)
		}

		if err != nil {
			return handleResponseError(err)
		}
		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualAPIVersions(args []string) error {
	apiID := args[0]
	ids := args[1:]

	r := make(resources.APIVersionSlice, len(ids))
	for i, id := range ids {
		resource, err := service.APIVersion(context.Background(), apiID, id)

		if err != nil {
			return handleResponseError(err)
		}

		r[i] = resource
	}

	printGetOutput(r)

	return nil
}
func getIndividualWorkspaces(args []string) error {
	r := make([]map[string]interface{}, len(args))
	uuidmap := prepareMap(resources.WorkspaceType)
	for i, name := range args {
		id, ok := uuidmap[name]
		var resource *resources.Workspace
		var err error
		if !ok {
			resource, err = service.Workspace(context.Background(), name)
		} else {
			resource, err = service.Workspace(context.Background(), id)
		}

		if err != nil {
			return handleResponseError(err)
		}

		data, err := json.Marshal(resource)
		if err != nil {
			return handleResponseError(err)
		}
		var tmp map[string]interface{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return handleResponseError(err)
		}
		keymap := make(map[string]int)
		for _, v := range strings.Split(ignoreKey.value, ",") {
			keymap[v] = 1
		}
		r[i] = util.ReformatMap(tmp, true, keymap)
	}

	printGetOutput(r)

	return nil
}

func getIndividualUser(args []string) error {
	resource, err := service.User(context.Background())

	if err != nil {
		return handleResponseError(err)
	}

	printGetOutput(resource)

	return nil
}

func getIndividualSchema(args []string) error {
	apiID := args[0]
	apiVersionID := args[1]
	id := args[2]

	resource, err := service.Schema(context.Background(), apiID, apiVersionID, id)

	if err != nil {
		return handleResponseError(err)
	}

	printGetOutput(resource)

	return nil
}

func getAPIRelations(apiID, apiVersionID string) error {
	resource, err := service.APIRelations(context.Background(), apiID, apiVersionID)

	if err != nil {
		return handleResponseError(err)
	}

	printGetOutput(resource)

	return nil
}

func getFormattedAPIRelations(apiID, apiVersionID string) error {
	resource, err := service.FormattedAPIRelationItems(context.Background(), apiID, apiVersionID)

	if err != nil {
		return handleResponseError(err)
	}

	printGetOutput(resource)

	return nil
}
