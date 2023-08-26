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

package registry

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/mod/semver"
	"sort"
	"strings"
)

func SkipPlatform(tag string) bool {
	if strings.HasSuffix(tag, "arm64") || strings.HasSuffix(tag, "amd64") {
		return false
	}
	return true
}

func ListTags(ctx context.Context, src string, filter func(tag string) bool) ([]string, error) {
	options := make([]crane.Option, 0)
	o := crane.GetOptions(options...)

	repo, err := name.NewRepository(src, o.Name...)
	if err != nil {
		return nil, fmt.Errorf("parsing repo %q: %w", src, err)
	}

	puller, err := remote.NewPuller(o.Remote...)
	if err != nil {
		return nil, err
	}

	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("reading tags for %s: %w", repo, err)
	}

	versions := make([]string, 0)

	for lister.HasNext() {
		tags, err := lister.Next(ctx)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags.Tags {
			newTag := tag
			if filter != nil && !filter(tag) {
				continue
			}
			versions = append(versions, newTag)
		}
	}
	semver.Sort(versions)
	sort.Sort(sort.Reverse(semver.ByVersion(versions)))
	return versions, nil
}
