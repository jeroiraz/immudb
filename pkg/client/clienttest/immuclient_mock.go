/*
Copyright 2019-2020 vChain, Inc.

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
package clienttest

import (
	"context"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc"
)

// ImmuClientMock ...
type ImmuClientMock struct {
	immuclient.ImmuClient

	SafeGetF       func(context.Context, []byte, ...grpc.CallOption) (*client.VerifiedItem, error)
	SafeSetF       func(context.Context, []byte, []byte) (*client.VerifiedIndex, error)
	SetF           func(context.Context, []byte, []byte) (*schema.Index, error)
	SafeReferenceF func(context.Context, []byte, []byte) (*client.VerifiedIndex, error)
	SafeZAddF      func(context.Context, []byte, float64, []byte) (*client.VerifiedIndex, error)
	HistoryF       func(context.Context, []byte) (*schema.StructuredItemList, error)
}

// SafeGet ...
func (icm *ImmuClientMock) SafeGet(ctx context.Context, key []byte, opts ...grpc.CallOption) (*client.VerifiedItem, error) {
	return icm.SafeGetF(ctx, key)
}

// SafeSet ...
func (icm *ImmuClientMock) SafeSet(ctx context.Context, key []byte, value []byte) (*client.VerifiedIndex, error) {
	return icm.SafeSetF(ctx, key, value)
}

// Set ...
func (icm *ImmuClientMock) Set(ctx context.Context, key []byte, value []byte) (*schema.Index, error) {
	return icm.SetF(ctx, key, value)
}

// SafeReference ...
func (icm *ImmuClientMock) SafeReference(ctx context.Context, reference []byte, key []byte) (*client.VerifiedIndex, error) {
	return icm.SafeReferenceF(ctx, reference, key)
}

// SafeZAdd ...
func (icm *ImmuClientMock) SafeZAdd(ctx context.Context, set []byte, score float64, key []byte) (*client.VerifiedIndex, error) {
	return icm.SafeZAddF(ctx, set, score, key)
}

// History ...
func (icm *ImmuClientMock) History(ctx context.Context, key []byte) (*schema.StructuredItemList, error) {
	return icm.HistoryF(ctx, key)
}
