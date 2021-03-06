// Copyright 2017 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encryption_test

import (
	"testing"

	"github.com/palantir/go-encrypted-config-value/encryption"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAESKeySerDe(t *testing.T) {
	aesKey, err := encryption.NewAESKey(256)
	require.NoError(t, err)

	aesKeyBytes := aesKey.Bytes()
	aesKey = encryption.AESKeyFromBytes(aesKeyBytes)

	cipher := encryption.NewAESGCMCipher()
	plaintext := "input plaintext"
	encrypted, err := cipher.Encrypt([]byte(plaintext), aesKey)
	require.NoError(t, err)

	decrypted, err := cipher.Decrypt(encrypted, aesKey)
	require.NoError(t, err)

	assert.Equal(t, plaintext, string(decrypted))
}
