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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringArgs(t *testing.T) {
	testCases := []struct {
		args      []string
		start     int
		end       int
		expectErr bool
		expected  string
	}{
		{[]string{"foo"}, 1, 2, false, "foo"},
		{[]string{"bar"}, 3, 4, false, "bar"},
		{[]string{}, 5, 6, true, ""},
		{[]string{"foo", "bar"}, 7, 8, true, ""},
		{[]string{""}, 9, 10, true, ""},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		node, err := parseStringArgs(tc.args, tc.start, tc.end)
		if tc.expectErr {
			assert.Error(err)
			assert.Nil(node)
			continue
		}

		assert.NoError(err)
		assert.NotNil(node)
		assert.Equal(tc.expected, node.GetValue())
		assert.Equal(tc.start, node.GetStartLine())
		assert.Equal(tc.end, node.GetEndLine())
	}
}
