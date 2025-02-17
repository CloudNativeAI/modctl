/*
 *     Copyright 2024 The CNAI Authors
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

package processor

import (
	"context"

	"github.com/CloudNativeAI/modctl/pkg/storage"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// NewReadmeProcessor creates a new README processor.
func NewReadmeProcessor(store storage.Storage, mediaType string, patterns []string) Processor {
	return &readmeProcessor{
		base: &base{
			store:     store,
			mediaType: mediaType,
			patterns:  patterns,
		},
	}
}

// readmeProcessor is the processor to process the README file.
type readmeProcessor struct {
	base *base
}

func (p *readmeProcessor) Name() string {
	return "readme"
}

func (p *readmeProcessor) Process(ctx context.Context, workDir, repo string) ([]ocispec.Descriptor, error) {
	return p.base.Process(ctx, workDir, repo)
}
