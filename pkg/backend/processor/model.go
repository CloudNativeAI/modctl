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
	"os"
	"regexp"

	"github.com/CloudNativeAI/modctl/pkg/backend/build"
	"github.com/CloudNativeAI/modctl/pkg/storage"
	modelspec "github.com/CloudNativeAI/model-spec/specs-go/v1"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// NewModelProcessor creates a new model processor.
func NewModelProcessor(models []string) Processor {
	return &modelProcessor{
		models: models,
	}
}

// modelProcessor is the processor to process the model file.
type modelProcessor struct {
	// models is the list of regular expressions to match the model file.
	models []string
}

func (p *modelProcessor) Name() string {
	return "model"
}

func (p *modelProcessor) Identify(_ context.Context, path string, info os.FileInfo) bool {
	for _, model := range p.models {
		if matched, _ := regexp.MatchString(model, info.Name()); matched {
			return true
		}
	}

	return false
}

func (p *modelProcessor) Process(ctx context.Context, store storage.Storage, repo, path, workDir string) (ocispec.Descriptor, error) {
	desc, err := build.BuildLayer(ctx, store, repo, path, workDir)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	// add model annotations.
	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}

	desc.Annotations[modelspec.AnnotationModel] = "true"
	return desc, nil
}
