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

package markdown

import (
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/cluster-image-docs/generator/types"
	"html/template"
	"os"
	"path"
	"strings"
)

const otherTemplate = `# {{.Type}} Image Versions
{{ $registry := .Registry }}{{ $repo := .Repo }}
Here are the versions of the images along with their corresponding links:
{{range .Data}}{{ $data := . }}
### [{{ $data.Name }}]({{$data.Url}})

{{range .Tags}}- [{{$registry}}/{{ $repo }}/{{$data.Name }}:{{ . }}](https://explore.ggcr.dev/?image={{$registry}}/{{$repo}}/{{ $data.Name }}:{{ . }})
{{end}}
{{end}}
`

type Type string

const (
	Rootfs Type = "Rootfs"
	Sealos Type = "Sealos"
	Laf    Type = "Laf"
	Apps   Type = "Apps"
)

type other struct {
	Type     Type
	Data     []types.ImageInfo
	Registry string
	Repo     string
}

func (r *other) Generator(dir string) error {
	fname := r.filename()
	f, err := os.Create(path.Join(dir, fname))
	if err != nil {
		return err
	}
	defer f.Close()
	t := template.Must(template.New("markdown").Parse(otherTemplate))
	err = t.Execute(f, r)
	if err != nil {
		return err
	}
	logger.Info("generator markdown %s success", fname)
	return nil
}

func (r *other) filename() string {
	switch r.Type {
	case Rootfs:
		return "rootfs.md"
	case Sealos:
		return "sealos.md"
	case Laf:
		return "laf.md"
	case Apps:
		return "apps.md"
	}
	return ""
}

func New(mdType Type, repoAddr string, data []types.ImageInfo) Template {
	registryArr := strings.Split(repoAddr, "/")
	return &other{
		Type:     mdType,
		Data:     data,
		Registry: registryArr[0],
		Repo:     registryArr[1],
	}
}
