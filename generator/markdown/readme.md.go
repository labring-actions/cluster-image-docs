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
	"fmt"
	"os"
	"path"
	"strings"
)

const readmeTemplate = `# Cluster Image Directory

Welcome to our documentation folder, where we list four distinct types of cluster images: Rootfs Images, Sealos-related Images, Laf-related Images, and Images for various applications.

- [Rootfs Images](./rootfs.md)
- [Sealos Images](./sealos.md)
- [Laf Images](./laf.md)
- [Application Images](./apps.md)

Please navigate to the above links to access the detailed documentation for each image category.

## Image Repository

The images in this documentation are hosted on our dedicated image repository:  %s

`

type readme struct {
	addr string
}

func (r *readme) Generator(dir string) error {
	return os.WriteFile(path.Join(dir, "README.md"), []byte(fmt.Sprintf(readmeTemplate, "`"+r.addr+"`")), 0644)
}

func NewReadme(addr string) Template {
	registryArr := strings.Split(addr, "/")
	return &readme{addr: registryArr[0]}
}
