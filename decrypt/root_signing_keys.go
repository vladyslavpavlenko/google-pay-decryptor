// Copyright (c) 2022 Rakhat

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package decrypt

import (
	"encoding/json"

	"github.com/vladyslavpavlenko/google-pay-decryptor/decrypt/types"
)

type RootSigningKey struct{}

func (r *RootSigningKey) Filter(rootKeys []byte) (types.RootKeys, []string, error) {
	filteredRootKeys, keyValues, err := loadRootSigningKeys(rootKeys)
	if err != nil {
		return types.RootKeys{}, nil, err
	}
	return filteredRootKeys, keyValues, nil
}

func loadRootSigningKeys(rootKeys []byte) (types.RootKeys, []string, error) {
	var keys types.RootSigningKey
	json.Unmarshal(rootKeys, &keys)
	keyValues := make([]string, 0)
	for _, filtered := range keys.RootKeys {
		keyValues = append(keyValues, filtered.KeyValue)
	}
	for _, filtered := range keys.RootKeys {
		if filtered.ProtocolVersion == "ECv2" {
			return filtered, keyValues, nil
		}
	}
	return types.RootKeys{}, keyValues, ErrParseJson
}
