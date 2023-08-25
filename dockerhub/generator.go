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

func generatorSyncFile(dir, key string, repos RepoInfo) error {
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
