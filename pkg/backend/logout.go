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

package backend

import (
	"context"

	"oras.land/oras-go/v2/registry/remote/credentials"
)

// Logout logs out of a registry.
func (b *backend) Logout(ctx context.Context, registry string) error {
	// read credentials from docker store.
	store, err := credentials.NewStoreFromDocker(credentials.StoreOptions{AllowPlaintextPut: true})
	if err != nil {
		return err
	}

	// remove credentials from store.
	if err := credentials.Logout(ctx, store, registry); err != nil {
		return err
	}

	return nil
}
