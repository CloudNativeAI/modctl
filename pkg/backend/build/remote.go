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

package build

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/CloudNativeAI/modctl/pkg/archiver"

	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v1"
	sha256 "github.com/minio/sha256-simd"
	godigest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func NewRemoteOutput(cfg *config, repo, tag string) (OutputStrategy, error) {
	remote, err := remote.NewRepository(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote repository: %w", err)
	}

	// gets the credentials store.
	credStore, err := credentials.NewStoreFromDocker(credentials.StoreOptions{AllowPlaintextPut: true})
	if err != nil {
		return nil, fmt.Errorf("failed to create credential store: %w", err)
	}

	httpClient := &http.Client{
		Transport: retry.NewTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.insecure,
			},
		}),
	}
	remote.Client = &auth.Client{
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
		Client:     httpClient,
	}

	remote.PlainHTTP = cfg.plainHTTP

	return &remoteOutput{
		cfg:    cfg,
		remote: remote,
		repo:   repo,
		tag:    tag,
	}, nil
}

type remoteOutput struct {
	cfg    *config
	remote *remote.Repository
	repo   string
	tag    string
}

// OutputLayer outputs the layer blob to the remote storage.
func (ro *remoteOutput) OutputLayer(ctx context.Context, mediaType, workDir, relPath string, reader io.Reader) (ocispec.Descriptor, error) {
	hash := sha256.New()
	size, err := io.Copy(hash, reader)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to copy layer to hash: %w", err)
	}

	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    godigest.Digest(fmt.Sprintf("sha256:%x", hash.Sum(nil))),
		Size:      size,
		Annotations: map[string]string{
			modelspec.AnnotationFilepath: relPath,
		},
	}

	reader, err = archiver.Tar(filepath.Join(workDir, relPath), workDir)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to create tar archive: %w", err)
	}

	exist, err := ro.remote.Blobs().Exists(ctx, desc)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to check if blob exists: %w", err)
	}

	if exist {
		return desc, nil
	}

	if err = ro.remote.Blobs().Push(ctx, desc, reader); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push layer to storage: %w", err)
	}

	return desc, nil
}

// OutputConfig outputs the config blob to the remote storage.
func (ro *remoteOutput) OutputConfig(ctx context.Context, mediaType string, configJSON []byte) (ocispec.Descriptor, error) {
	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    godigest.Digest(fmt.Sprintf("sha256:%x", sha256.Sum256(configJSON))),
		Size:      int64(len(configJSON)),
	}

	exist, err := ro.remote.Blobs().Exists(ctx, desc)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to check if blob exists: %w", err)
	}

	if exist {
		return desc, nil
	}

	if err = ro.remote.Blobs().Push(ctx, desc, bytes.NewReader(configJSON)); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push config to storage: %w", err)
	}

	return desc, nil
}

// OutputManifest outputs the manifest blob to the remote storage.
func (ro *remoteOutput) OutputManifest(ctx context.Context, mediaType string, manifestJSON []byte) (ocispec.Descriptor, error) {
	desc := ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    godigest.Digest(fmt.Sprintf("sha256:%x", sha256.Sum256(manifestJSON))),
		Size:      int64(len(manifestJSON)),
	}

	exist, err := ro.remote.Manifests().Exists(ctx, desc)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to check if blob exists: %w", err)
	}

	if exist {
		return desc, nil
	}

	if err = ro.remote.Manifests().Push(ctx, desc, bytes.NewReader(manifestJSON)); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push manifest to storage: %w", err)
	}

	// Tag the manifest.
	if err = ro.remote.Tag(ctx, desc, ro.tag); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to tag manifest: %w", err)
	}

	return desc, nil
}
