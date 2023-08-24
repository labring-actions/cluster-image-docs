/*
Copyright 2023 cuisongliu@qq.com.

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

package dockerhub

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"gopkg.in/yaml.v2"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const tmpl = `docker.io:
  {{- if .ByTagRegex }}
  images-by-tag-regex:
    {{- range .Repos }}
    labring/{{ .Name }}: {{ .Filter }}
    {{- end }}
  {{- else }}
  images:
    {{- range .Repos }}
    labring/{{ .Name }}: []
    {{- end }}
  {{- end }}
  tls-verify: false
`

const workflowTmpl = `name: {{ .SYNC_FILE_NAME }} 
on:
  push:
    branches: [ main ]
    paths:
      - "{{ .SYNC_FILE }}"
      - ".github/workflows/{{ .SYNC_FILE_NAME }}"
  schedule:
    - cron: '0 16 * * *'
  workflow_dispatch:

env:
  USERNAME: {{ .USER_KEY }}
  PASSWORD: {{ .PASSWORD_KEY }}

jobs:
  image-sync:
    runs-on: {{.RUN_ON}}

    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: check podman
        run: |
          sudo podman version

      - name: sync images
        run: |
          sudo podman run -it --rm -v ${PWD}:/workspace -w /workspace quay.io/skopeo/stable:latest \
          sync --src yaml --dest docker {{ .SYNC_FILE }} {{ .REGISTRY_KEY }}/{{ .REPOSITORY_KEY }} \
          --dest-username $USERNAME --dest-password "$PASSWORD" \
          --keep-going --retry-times 2 --all
`

const prefix = "auto-sync"

func autoRemoveGenerator(dir string) error {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果不是目录，并且文件名以指定的前缀开始，就删除它
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			os.Remove(path)
		}

		return nil
	}); err != nil {
		return err
	}
	logger.Info("auto remove %s files success", dir)
	return nil
}

func generatorSyncFile(dir, key string, repos RepoInfoList) error {
	f, err := os.Create(path.Join(dir, fmt.Sprintf("%s-%s.yaml", prefix, key)))
	if err != nil {
		return err
	}
	defer f.Close()
	t := template.Must(template.New("repos").Parse(tmpl))
	err = t.Execute(f, repos)
	if err != nil {
		return err
	}
	logger.Info("generator sync config %s-%s.yaml success", prefix, key)
	return nil
}

func getCIRun(file string) ([]string, error) {
	type Runner struct {
		Name         string   `yaml:"name"`
		Cloud        string   `yaml:"cloud"`
		MachineImage string   `yaml:"machine_image"`
		InstanceType []string `yaml:"instance_type"`
		Labels       []string `yaml:"labels"`
		Region       []string `yaml:"region"`
	}

	type Config struct {
		Runners []Runner `yaml:"runners"`
	}

	data, err := os.ReadFile(file) //replace with your config file
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	if len(config.Runners) == 0 {
		return nil, fmt.Errorf("not found runners")
	}
	if len(config.Runners[0].Labels) == 0 {
		return nil, fmt.Errorf("not found labels")
	}
	for _, label := range config.Runners[0].Labels {
		if !strings.HasPrefix(label, "cirun") {
			return nil, fmt.Errorf("not found cirun labels, must has prefix %s", "cirun")
		}
	}
	return config.Runners[0].Labels, nil
}

func generatorWorkflowFile(dir, syncDir, key string, labels []string) error {
	syncFileName := fmt.Sprintf("%s-%s.yaml", prefix, key)
	syncFile := path.Join(syncDir, syncFileName)
	f, err := os.Create(path.Join(dir, fmt.Sprintf("%s-%s.yaml", prefix, key)))
	if err != nil {
		return err
	}
	defer f.Close()
	t := template.Must(template.New("repos").Parse(workflowTmpl))

	runOn := "ubuntu-22.04"
	if len(labels) > 0 {
		runOn = fmt.Sprintf("%s--${{ github.run_id }}", labels[0])
	}

	err = t.Execute(f, map[string]string{
		"PREFIX":         prefix,
		"SYNC_FILE":      syncFile,
		"SYNC_FILE_NAME": syncFileName,
		"RUN_ON":         runOn,
		"USER_KEY":       "${{ vars.A_REGISTRY_USERNAME }}",
		"PASSWORD_KEY":   "${{ secrets.A_REGISTRY_TOKEN }}",
		"REGISTRY_KEY":   "${{ vars.A_REGISTRY_NAME }}",
		"REPOSITORY_KEY": "${{ vars.A_REGISTRY_REPOSITORY }}",
	})
	if err != nil {
		return err
	}
	logger.Info("generator workflow config %s-%s.yaml success", prefix, key)
	return nil
}
