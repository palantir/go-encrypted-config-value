// Copyright 2017 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encryptedconfigvalue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/go-encrypted-config-value/encryptedconfigvalue"
)

type aliasForStringType string

type doubleAliasForStringType aliasForStringType

type complicatedStruct struct {
	String        string
	Pointer       *string
	DoublePointer **string
	Slice         []string
	Array         [2]string
	Map           map[string]string
	Alias         aliasForStringType
	DoubleAlias   doubleAliasForStringType
	Interface     interface{}
	Nested        map[string][]interface{}
	unexported    string
}

type decryptVarsTestCase struct {
	name string
	in   interface{}
	want interface{}
}

const (
	encryptedValVar = "${enc:eyJ0eXBlIjoiQUVTIiwibW9kZSI6IkdDTSIsImNpcGhlcnRleHQiOiJNOTRrSXlvYTUrMloiLCJpdiI6InVBR3FSbFA5d2l6cGRCMHoiLCJ0YWciOiJBQ1N1ekR3VFVMb21zanhwRk1rWUtBPT0ifQ==}"
	decrypted       = "plaintext"
)

var aesKeyWithType = encryptedconfigvalue.MustNewKeyWithType("AES:LICx0yKzQm5a6IE13aJ3xOsRv+8AujqHocTFI4yk4Jw=")

func decryptVarsTestCases() []decryptVarsTestCase {
	return []decryptVarsTestCase{
		{
			"Single string input",
			encryptedValVar,
			decrypted,
		},
		{
			"Multiple string input",
			"hello ${enc:eyJ0eXBlIjoiQUVTIiwibW9kZSI6IkdDTSIsImNpcGhlcnRleHQiOiJNOTRrSXlvYTUrMloiLCJpdiI6InVBR3FSbFA5d2l6cGRCMHoiLCJ0YWciOiJBQ1N1ekR3VFVMb21zanhwRk1rWUtBPT0ifQ==} this is ${enc:eyJ0eXBlIjoiQUVTIiwibW9kZSI6IkdDTSIsImNpcGhlcnRleHQiOiJNOTRrSXlvYTUrMloiLCJpdiI6InVBR3FSbFA5d2l6cGRCMHoiLCJ0YWciOiJBQ1N1ekR3VFVMb21zanhwRk1rWUtBPT0ifQ==} ",
			"hello plaintext this is plaintext ",
		},
		{
			"Indirect string input",
			testStringPtr(encryptedValVar),
			testStringPtr(decrypted),
		},
		{
			"Complicated struct",
			complicatedStruct{
				String:        encryptedValVar,
				Pointer:       testStringPtr(encryptedValVar),
				DoublePointer: testStringDoublePtr(encryptedValVar),
				Slice: []string{
					encryptedValVar,
					"hello",
				},
				Array: [2]string{
					"goodbye",
					encryptedValVar,
				},
				Map: map[string]string{
					"key":           encryptedValVar,
					encryptedValVar: "value",
				},
				Alias:       aliasForStringType(encryptedValVar),
				DoubleAlias: doubleAliasForStringType(encryptedValVar),
				Interface:   encryptedValVar,
				Nested: map[string][]interface{}{
					"key": {
						complicatedStruct{
							String: encryptedValVar,
						},
					},
				},
				unexported: encryptedValVar,
			},
			complicatedStruct{
				String:        decrypted,
				Pointer:       testStringPtr(decrypted),
				DoublePointer: testStringDoublePtr(decrypted),
				Slice: []string{
					decrypted,
					"hello",
				},
				Array: [2]string{
					"goodbye",
					decrypted,
				},
				Map: map[string]string{
					"key":           decrypted,
					encryptedValVar: "value",
				},
				Alias:       aliasForStringType(decrypted),
				DoubleAlias: doubleAliasForStringType(decrypted),
				Interface:   decrypted,
				Nested: map[string][]interface{}{
					"key": {
						complicatedStruct{
							String: decrypted,
						},
					},
				},
				unexported: encryptedValVar,
			},
		},
	}
}

func TestDecryptEncryptedStringVariables(t *testing.T) {
	for i, currCase := range decryptVarsTestCases() {
		got := currCase.in
		// pass in address of "got" to modify in-place
		encryptedconfigvalue.DecryptEncryptedStringVariables(&got, aesKeyWithType)
		assert.Equal(t, currCase.want, got, "Case %d: %s", i, currCase.name)
	}
}

func TestCopyWithEncryptedStringVariablesDecrypted(t *testing.T) {
	for i, currCase := range decryptVarsTestCases() {
		// pass "currCase.in" directly
		got := encryptedconfigvalue.CopyWithEncryptedStringVariablesDecrypted(currCase.in, aesKeyWithType)
		assert.Equal(t, currCase.want, got, "Case %d: %s", i, currCase.name)
	}
}

// TestCopyWithEncryptedStringVariablesDecryptedPerformsShallowCopy verifies that the "copy" performed by
// CopyWithEncryptedStringVariablesDecrypted is a shallow copy of the input (because it is performed in a manner that is
// equivalent to declaring a new variable of the input type and performing an assignment). Thus, if the input is a
// pointer, the "copy" and the original will both point to the same value and thus the modifications will impact both.
func TestCopyWithEncryptedStringVariablesDecryptedPerformsShallowCopy(t *testing.T) {
	// CopyWithEncryptedStringVariablesDecrypted performs a "copy" of the provided input, but the copy is equivalent
	// to declaring a new variable and performing an assignment (shallow copy).
	type inStruct struct {
		Val string
	}

	in := &inStruct{
		Val: encryptedValVar,
	}
	rawGot := encryptedconfigvalue.CopyWithEncryptedStringVariablesDecrypted(in, aesKeyWithType)

	got, ok := rawGot.(*inStruct)
	require.True(t, ok, "type of returned value was not correct")
	assert.Equal(t, decrypted, got.Val, "value should be decrypted")
	assert.Equal(t, in, got, "values should be the same because input was a pointer")
	assert.True(t, in == got, "pointer values should be the same because they were copied")
}

func testStringPtr(input string) *string {
	return &input
}

func testStringDoublePtr(input string) **string {
	ptr := testStringPtr(input)
	return &ptr
}
