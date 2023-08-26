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

package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/cluster-image-docs/generator/registry"
	"github.com/labring-actions/cluster-image-docs/generator/types"
	"github.com/labring-actions/cluster-image-docs/generator/utils"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

var specialRepos = []string{"kubernetes", "kubernetes-crio", "kubernetes-docker"}

func fetchDockerHubAllRepo(registryRepoName string) (*types.RepoInfo, error) {
	type Repo struct {
		Name string `json:"name"`
	}

	type Repositories struct {
		Results []Repo `json:"results"`
		Next    string `json:"next"`
	}
	dockerAPI := "https://hub.docker.com/v2/repositories/labring?page_size=10"
	repos := &types.RepoInfo{
		Rootfs: make([]types.ImageInfo, 0),
		Sealos: make([]types.ImageInfo, 0),
		Laf:    make([]types.ImageInfo, 0),
		Apps:   make([]types.ImageInfo, 0),
	}
	var reposMutex sync.Mutex

	if err := utils.Retry(func() error {
		for dockerAPI != "" {
			logger.Debug("fetch dockerhub url: %s", dockerAPI)
			data, err := utils.Request(dockerAPI, "GET", []byte(""), 0)
			if err != nil {
				return err
			}
			var repositories Repositories
			if err = json.Unmarshal(data, &repositories); err != nil {
				return err
			}
			eg, _ := errgroup.WithContext(context.Background())

			for _, repo := range repositories.Results {
				repoName := repo.Name
				eg.Go(func() error {
					filterFn := registry.SkipPlatform
					if strings.HasPrefix(repoName, "sealos") || strings.HasPrefix(repoName, "laf") {
						filterFn = func(tag string) bool {
							if strings.HasPrefix(tag, "v") || tag == "latest" {
								return registry.SkipPlatform(tag)
							}
							return false
						}
					}

					tags, err := registry.ListTags(context.Background(), fmt.Sprintf("%s/%s", registryRepoName, repoName), filterFn)
					if err != nil {
						logger.Warn("ListTags %s/%s error: %s", registryRepoName, repoName, err.Error())
					}
					types.Sort(tags)
					reposMutex.Lock()
					defer reposMutex.Unlock()
					info := types.ImageInfo{
						Name: repoName,
						Tags: tags,
					}

					if utils.StringInSlice(repoName, specialRepos) {
						repos.Rootfs = append(repos.Rootfs, info)
					} else if strings.HasPrefix(repoName, "sealos") {
						repos.Sealos = append(repos.Sealos, info)
					} else if strings.HasPrefix(repoName, "laf") {
						repos.Laf = append(repos.Laf, info)
					} else {
						repos.Apps = append(repos.Apps, info)
					}
					return nil
				})
			}
			if err = eg.Wait(); err != nil {
				return err
			}
			dockerAPI = repositories.Next
		}
		return nil
	}); err != nil {
		logger.Error("get dockerhub repo error: %s", err.Error())
		return nil, err
	}

	return repos, nil
}
