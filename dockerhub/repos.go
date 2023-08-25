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
	"encoding/json"
	"github.com/cuisongliu/logger"
	"strings"
)

type RepoInfo struct {
	Rootfs []string
	Sealos []string
	Apps   []string
}

type RepoVersionMap map[string][]string

var specialRepos = []string{"kubernetes", "kubernetes-crio", "kubernetes-docker"}

func fetchDockerHubAllRepo() (*RepoInfo, error) {
	type Repo struct {
		Name string `json:"name"`
	}

	type Repositories struct {
		Results []Repo `json:"results"`
		Next    string `json:"next"`
	}
	fetchURL := "https://hub.docker.com/v2/repositories/labring?page_size=10"
	versions := &RepoInfo{
		Rootfs: make([]string, 0),
		Sealos: make([]string, 0),
		Apps:   make([]string, 0),
	}
	if err := Retry(func() error {
		for fetchURL != "" {
			logger.Debug("fetch dockerhub url: %s", fetchURL)
			data, err := Request(fetchURL, "GET", []byte(""), 0)
			if err != nil {
				return err
			}
			var repositories Repositories
			if err = json.Unmarshal(data, &repositories); err != nil {
				return err
			}
			for _, repo := range repositories.Results {
				if stringInSlice(repo.Name, specialRepos) {
					versions.Rootfs = append(versions.Rootfs, repo.Name)
				} else if strings.HasPrefix(repo.Name, "sealos") {
					if strings.HasPrefix(repo.Name, "sealos-cloud") || repo.Name == "sealos" || repo.Name == "sealos-patch" {
						versions.Sealos = append(versions.Sealos, repo.Name)
					}
				} else if strings.HasPrefix(repo.Name, "laf") {
					versions.Sealos = append(versions.Sealos, repo.Name)
				} else {
					versions.Apps = append(versions.Apps, repo.Name)
				}
			}
			fetchURL = repositories.Next
		}
		return nil
	}); err != nil {
		logger.Error("get dockerhub repo error: %s", err.Error())
		return nil, err
	}
	return versions, nil
}
