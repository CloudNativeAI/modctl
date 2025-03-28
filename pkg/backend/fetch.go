/*
 *     Copyright 2025 The CNAI Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"

	internalpb "github.com/CloudNativeAI/modctl/internal/pb"
	"github.com/CloudNativeAI/modctl/pkg/backend/remote"
	"github.com/CloudNativeAI/modctl/pkg/config"
)

// Fetch fetches partial files to the output.
func (b *backend) Fetch(ctx context.Context, target string, cfg *config.Fetch) error {
	// parse the repository and tag from the target.
	ref, err := ParseReference(target)
	if err != nil {
		return fmt.Errorf("failed to parse the target: %w", err)
	}

	repo, tag := ref.Repository(), ref.Tag()
	client, err := remote.New(repo, remote.WithPlainHTTP(cfg.PlainHTTP), remote.WithInsecure(cfg.Insecure))
	if err != nil {
		return fmt.Errorf("failed to create remote client: %w", err)
	}

	_, manifestReader, err := client.Manifests().FetchReference(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to fetch the manifest: %w", err)
	}

	defer manifestReader.Close()

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manifestReader).Decode(&manifest); err != nil {
		return fmt.Errorf("failed to decode the manifest: %w", err)
	}

	layers := []ocispec.Descriptor{}
	// filter the layers by patterns.
	for _, layer := range manifest.Layers {
		for _, pattern := range cfg.Patterns {
			if anno := layer.Annotations; anno != nil {
				matched, err := filepath.Match(pattern, anno[modelspec.AnnotationFilepath])
				if err != nil {
					return fmt.Errorf("failed to match pattern: %w", err)
				}

				if matched {
					layers = append(layers, layer)
				}
			}
		}
	}

	if len(layers) == 0 {
		return fmt.Errorf("no layers matched the patterns")
	}

	pb := internalpb.NewProgressBar()
	pb.Start()
	defer pb.Stop()

	g := &errgroup.Group{}
	g.SetLimit(cfg.Concurrency)

	for _, layer := range layers {
		g.Go(func() error {
			return pullAndExtractFromRemote(ctx, pb, internalpb.NormalizePrompt("Fetching blob"), client, cfg.Output, layer)
		})
	}

	return g.Wait()
}
