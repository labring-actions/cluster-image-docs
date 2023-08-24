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
	"fmt"
	"github.com/cuisongliu/logger"
	"strings"
)

type RepoInfo struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
	Filter   string   `json:"filter"`
}

type RepoInfoList struct {
	Repos      []RepoInfo `json:"repos"`
	ByTagRegex bool
}

var specialRepos = []string{"kubernetes", "kubernetes-crio", "kubernetes-docker"}

func fetchDockerHubAllRepo() (map[string]RepoInfoList, error) {
	type Repo struct {
		Name string `json:"name"`
	}

	type Repositories struct {
		Results []Repo `json:"results"`
		Next    string `json:"next"`
	}

	fetchURL := "https://hub.docker.com/v2/repositories/labring?page_size=10"

	versions := make(map[string]RepoInfoList)
	if err := Retry(func() error {
		index := 0
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
			newRepos := make([]RepoInfo, 0)
			for _, repo := range repositories.Results {
				if stringInSlice(repo.Name, specialRepos) {
					versions[repo.Name] = RepoInfoList{
						Repos: []RepoInfo{
							{Name: repo.Name, Filter: "^v(1\\.2[0-9]\\.[1-9]?[0-9]?)(\\.)?$"},
						},
						ByTagRegex: true,
					}
				} else if strings.HasPrefix(repo.Name, "sealos") {
					if strings.HasPrefix(repo.Name, "sealos-cloud") || repo.Name == "sealos" || repo.Name == "sealos-patch" {
						versions[repo.Name] = RepoInfoList{
							Repos: []RepoInfo{
								{Name: repo.Name, Filter: "^v.*"},
							},
							ByTagRegex: true,
						}
					}
					//logger.Warn("sealos container image repo is deprecated, please use sealos cloud repo")
				} else if strings.HasPrefix(repo.Name, "laf") {
					versions[repo.Name] = RepoInfoList{
						Repos: []RepoInfo{
							{Name: repo.Name, Filter: "^v.*"},
						},
						ByTagRegex: true,
					}
				} else {
					newRepos = append(newRepos, RepoInfo{Name: repo.Name})
				}
			}
			versions[fmt.Sprintf("image-%d", index)] = RepoInfoList{
				Repos:      newRepos,
				ByTagRegex: false,
			}
			index++
			fetchURL = repositories.Next
		}
		return nil
	}); err != nil {
		logger.Error("get dockerhub repo error: %s", err.Error())
		return nil, err
	}
	return versions, nil
}
