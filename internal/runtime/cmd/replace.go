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
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevinswiber/postmanctl/pkg/sdk/resources"
	"github.com/spf13/cobra"
)

var mode string

var diffFile defaultValue

func init() {
	replaceCmd := &cobra.Command{
		Use:   "replace",
		Short: "Replace existing Postman resources.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				inputReader = os.Stdin
			} else {
				if inputFile == "" {
					return errors.New("flag \"filename\" not set, use \"--filename\" or stdin")
				}
			}

			return nil
		},
	}
	replaceCmd.PersistentFlags().StringVarP(&inputFile, "filename", "f", "", "the filename used to replace the resource (required when not using data from stdin)")
	replaceCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "force", "force/compare -> force replace or compare before replace")
	//replaceCmd.PersistentFlags().StringVarP(&diffFile, "diff", "df", "", "the file diff report. Default diff report will print to console")
	replaceCmd.PersistentFlags().VarP(&ignoreKey, "ignore-key", "i", "ignore json key in response")
	replaceCmd.PersistentFlags().VarP(&diffFile, "diff-file", "d", "ignore json key in response")

	replaceCmd.AddCommand(
		generateReplaceSubcommand(resources.CollectionType, "collection", []string{"co"}),
		generateReplaceSubcommand(resources.EnvironmentType, "environment", []string{"env"}),
		generateReplaceSubcommand(resources.MonitorType, "monitor", []string{"mon"}),
		generateReplaceSubcommand(resources.MockType, "mock", []string{}),
		generateReplaceSubcommand(resources.WorkspaceType, "workspace", []string{"ws"}),
		generateReplaceSubcommand(resources.APIType, "api", []string{}),
		generateReplaceSubcommand(resources.APIVersionType, "api-version", []string{}),
		generateReplaceSubcommand(resources.SchemaType, "schema", []string{}),
	)

	rootCmd.AddCommand(replaceCmd)
}

func generateReplaceSubcommand(t resources.ResourceType, use string, aliases []string) *cobra.Command {
	cmd := cobra.Command{
		Use:     use,
		Aliases: aliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return replaceResource(t, args[0])
		},
	}

	if t == resources.APIVersionType || t == resources.SchemaType {
		cmd.Flags().StringVar(&forAPI, "for-api", "", "the associated API ID (required)")
		cmd.MarkFlagRequired("for-api")
	}

	if t == resources.SchemaType {
		cmd.Flags().StringVar(&forAPIVersion, "for-api-version", "", "the associated API Version ID (required)")
		cmd.MarkFlagRequired("for-api-version")
	}

	return &cmd
}

func replaceResource(t resources.ResourceType, resourceID string) error {
	if inputReader == nil {
		r, err := os.Open(inputFile)

		if err != nil {
			return err
		}

		defer r.Close()

		inputReader = r
	}

	var (
		id  string
		err error
	)
	uuid := prepareMap(t)
	tmpID, ok := uuid[resourceID]
	if ok {
		resourceID = tmpID
	}
	switch mode {
	case "force":
	case "compare":
		if len(inputFile) == 0 {
			fmt.Fprintf(os.Stderr, "compare mode only work with file")
			return nil
		}
		if !doCompare(inputFile, resourceID, t) {
			return nil
		}
	default:
		fmt.Fprintf(os.Stderr, "wrong mode %s, mode must be force or compare\n", mode)
		return nil
	}

	ctx := context.Background()
	switch t {
	case resources.CollectionType:
		id, err = service.ReplaceCollectionFromReader(ctx, inputReader, resourceID)
	case resources.EnvironmentType:
		id, err = service.ReplaceEnvironmentFromReader(ctx, inputReader, resourceID)
	case resources.MockType:
		id, err = service.ReplaceMockFromReader(ctx, inputReader, resourceID)
	case resources.MonitorType:
		id, err = service.ReplaceMonitorFromReader(ctx, inputReader, resourceID)
	case resources.WorkspaceType:
		id, err = service.ReplaceWorkspaceFromReader(ctx, inputReader, resourceID)
	case resources.APIType:
		id, err = service.ReplaceAPIFromReader(ctx, inputReader, resourceID)
	case resources.APIVersionType:
		id, err = service.ReplaceAPIVersionFromReader(ctx, inputReader, forAPI, resourceID)
	case resources.SchemaType:
		id, err = service.ReplaceSchemaFromReader(ctx, inputReader, resourceID, forAPI, forAPIVersion)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(id)

	return nil
}

func doCompare(file string, resourceID string, t resources.ResourceType) bool {
	r, err := os.Open(file)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}
	defer r.Close()
	var data interface{}
	ctx := context.Background()
	switch t {
	case resources.CollectionType:
		data, err = service.Collection(ctx, resourceID)
	case resources.EnvironmentType:
		data, err = service.Environment(ctx, resourceID)
	case resources.MockType:
		data, err = service.Mock(ctx, resourceID)
	case resources.MonitorType:
		data, err = service.Monitor(ctx, resourceID)
	case resources.WorkspaceType:
		data, err = service.Workspace(ctx, resourceID)
	default:
		fmt.Fprintln(os.Stderr, "no supported type")
		return true
	}
	btmp, err := json.Marshal(data)
	if err != nil {
		return false
	}
	var tmp map[string]interface{}
	err = json.Unmarshal(btmp, &tmp)
	if err != nil {
		return false
	}
	keymap := make(map[string]int)
	for _, v := range strings.Split(ignoreKey.value, ",") {
		keymap[v] = 1
	}
	tmp = util.ReformatMap(tmp, true, keymap)
	b, err := ioutil.ReadAll(r)

	if err != nil {
		return false
	}

	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return false
	}
	dff := util.CompareMap("", tmp, v)
	tt, _ := json.MarshalIndent(dff, "", "  ")
	if len(diffFile.value) > 0 {
		fmt.Printf("Write diff report to file %s\n", diffFile.value)
		ioutil.WriteFile(diffFile.value, tt, 0644)
	} else {
		fmt.Println("Diff report:")
		fmt.Println(string(tt))
	}
	fmt.Println("Please check carefully before confirm replace.")
	var confirm string
	fmt.Printf("Are you sure to merge(Y/N):")
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) == "y" {
		return true
	}
	fmt.Printf("Cancel!\n")
	return false
}
