go-encrypted-config-value
=========================
`go-encrypted-config-value` is a Go implementation of the [encrypted-config-value](https://github.com/palantir/encrypted-config-value)
library. It provides a simple mechanism to encrypt and decrypt values stored in configuration on disk.


Overview
--------
The `encryption` package provides primitives such as `Key` and `Cipher`. It provides implementations of an AES-GCM
cipher with configurable nonce and tag sizes and an RSA-OAEP cipher with configurable hash algorithms for OAEP and MDF1.
It also provides functions for creating and serializing AES and RSA keys.

The `encryptedconfigvalue` package provides functionality for creating, decrypting and serializing/deserializing
encrypted config values.


Usage
-----
AES:

```go
keyPair, err := encryptedconfigvalue.AES.GenerateKeyPair()
encryptedVal, err := encryptedconfigvalue.AES.Encrypter().Encrypt("secret text", keyPair.EncryptionKey)

serializedKey := keyPair.EncryptionKey.ToSerializable()
serializedValue, err := encryptedVal.ToSerializable()

rehydratedKey, err := encryptedconfigvalue.NewKeyWithType(serializedKey)
rehydratedValue, err := encryptedconfigvalue.NewEncryptedValue(serializedValue)

plaintext, err := rehydratedValue.Decrypt(rehydratedKey)
```

RSA:

```go
keyPair, err := encryptedconfigvalue.RSA.GenerateKeyPair()
encryptedVal, err := encryptedconfigvalue.RSA.Encrypter().Encrypt("secret text", keyPair.EncryptionKey)

serializedPublicKey := keyPair.EncryptionKey.ToSerializable()
serializedPrivateKey := keyPair.DecryptionKey.ToSerializable()
serializedValue, err := encryptedVal.ToSerializable()

rehydratedDecryptionKey, err := encryptedconfigvalue.NewKeyWithType(serializedKey)
rehydratedValue, err := encryptedconfigvalue.NewEncryptedValue(serializedValue)

plaintext, err := rehydratedValue.Decrypt(rehydratedDecryptionKey)
```

Values in Configuration:

* `encryptedconfigvalue.ContainsEncryptedConfigValueStringVars` returns true if the provided input contains any entries
  values of the form "${enc:...}"
* `encryptedconfigvalue.DecryptAllEncryptedValueStringVars` returns a version of the provided input where all string
  values of the form "${enc:...}" are replaced with the result of decrypting the values using the provided key
* `encryptedconfigvalue.DecryptEncryptedStringVariables` recusrively finds all occurrences of string values of the form
  "${enc:...}" in the exported fields of an object and replaces them with the result of decrypting the values using the
  provided key


Backwards Compatibility
-----------------------
This library supports reading encryption keys and encrypted values that are stored in the legacy format. The
`encryptedconfigvalue.NewKeyWithType` and `encryptedconfigvalue.NewEncryptedValue` will both accept valid values in the
legacy format.

This library can generate `EncryptedValue` objects that serialize using legacy encrypters that are provided as part of
the library.

AES: `encryptedVal, err := encryptedconfigvalue.LegacyAESGCMEncrypter().Encrypt("secret text", keyPair.EncryptionKey)`
RSA: `encryptedVal, err := encryptedconfigvalue.LegacyRSAOAEPEncrypter().Encrypt("secret text", keyPair.EncryptionKey)`

The `ToSerializable` function for these legacy values will return a string that encodes the values using the legacy format.

The format for AES keys did not change between the legacy and new format.

The format for RSA keys did change, but because the legacy format for RSA keys was not widely used and had deficiencies
(it was not able to differentiate between a public and private key based on the serialized format), writing such values
is not supported by this library.


License
-------
This project is made available under the [BSD 3-Clause License](https://opensource.org/licenses/BSD-3-Clause).
