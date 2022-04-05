/*
Copyright © 2022 Yevhen Lebid ealebed@gmail.com

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

package output

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func MarshalToJson(input interface{}) ([]byte, error) {
	pretty, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to json: %v", err)
	}
	return pretty, nil
}

func MarshalToYaml(input interface{}) ([]byte, error) {
	pretty, err := yaml.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to yaml: %v", err)
	}
	return pretty, nil
}

func JsonOutput(input interface{}) {
	res, err := MarshalToJson(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(res))
}

func YamlOutput(input interface{}) {
	res, err := MarshalToYaml(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(res))
}
