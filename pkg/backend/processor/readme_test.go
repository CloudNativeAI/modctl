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
	"testing"
	"testing/fstest"

	modelspec "github.com/CloudNativeAI/modctl/pkg/spec"
	"github.com/CloudNativeAI/modctl/test/mocks/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadmeProcessor_Name(t *testing.T) {
	p := NewReadmeProcessor()
	assert.Equal(t, "readme", p.Name())
}

func TestReadmeProcessor_Identify(t *testing.T) {
	p := NewReadmeProcessor()
	mockFS := fstest.MapFS{
		"README":    &fstest.MapFile{},
		"README.md": &fstest.MapFile{},
		"LICENSE":   &fstest.MapFile{},
	}
	info, err := mockFS.Stat("README")
	assert.NoError(t, err)
	assert.True(t, p.Identify(context.Background(), "README", info))

	info, err = mockFS.Stat("README.md")
	assert.NoError(t, err)
	assert.True(t, p.Identify(context.Background(), "README.md", info))

	info, err = mockFS.Stat("LICENSE")
	assert.NoError(t, err)
	assert.False(t, p.Identify(context.Background(), "LICENSE", info))
}

func TestReadmeProcessor_Process(t *testing.T) {
	p := NewReadmeProcessor()
	ctx := context.Background()
	mockStore := &storage.Storage{}
	repo := "test-repo"
	path := "README"
	mockFS := fstest.MapFS{
		"README": &fstest.MapFile{},
	}
	info, err := mockFS.Stat("README")
	assert.NoError(t, err)

	mockStore.On("PushBlob", ctx, repo, mock.Anything).Return("sha256:1234567890abcdef", int64(1024), nil)

	desc, err := p.Process(ctx, mockStore, repo, path, info)
	assert.NoError(t, err)
	assert.NotNil(t, desc)
	assert.Equal(t, "sha256:1234567890abcdef", desc.Digest.String())
	assert.Equal(t, int64(1024), desc.Size)
	assert.Equal(t, "true", desc.Annotations[modelspec.AnnotationReadme])
}