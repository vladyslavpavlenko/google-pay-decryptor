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
	subsig "github.com/google/tink/go/signature/subtle"
	"github.com/vladyslavpavlenko/google-pay-decryptor/decrypt/types"
)

func VerifySignature(token types.Token, keyValues []string, receipientId string) error {
	if token.ProtocolVersion != "ECv2" {
		return ErrProtocolV
	}

	if err := verifyIntermediateSigningKey(token, keyValues); err != nil {
		return err
	}

	signedKey, err := token.IntermediateSigningKey.UnmarshalSignedKey(token.IntermediateSigningKey.SignedKey)
	if err != nil {
		return err
	}

	if !CheckTime(signedKey.KeyExpiration) {
		return ErrValidateTimeKey
	}

	if err := VerifyMessageSignature(signedKey.KeyValue, token, receipientId); err != nil {
		return err
	}

	return nil
}

func verifyIntermediateSigningKey(token types.Token, keyValues []string) error {
	signatures := token.IntermediateSigningKey.Signatures
	signedKey := token.IntermediateSigningKey.SignedKey
	signedData := ConstructSignature(GooglePaySenderID, token.ProtocolVersion, signedKey)
	for _, key := range keyValues {
		var pk *PublicKey
		publicKey, err := pk.LoadPublicKey(key)
		if err != nil {
			return err
		}
		for _, signature := range signatures {
			signatureDecoded, err := Base64Decode(signature)
			if err != nil {
				return err
			}
			verifyer, err := subsig.NewECDSAVerifierFromPublicKey(GooglePaySHA256HashAlgorithm, GooglePayDEREncoding, publicKey)
			if err != nil {
				return err
			}
			if isErr := verifyer.Verify(signatureDecoded, signedData); isErr != nil {
				continue
			}
			return nil
		}
	}
	return ErrVerifySignature
}

func VerifyMessageSignature(keyValue string, token types.Token, receipientId string) error {
	var pk PublicKey
	publicKey, err := pk.LoadPublicKey(keyValue)
	if err != nil {
		return err
	}
	signature, _ := Base64Decode(token.Signature)
	signedData := ConstructSignature(GooglePaySenderID, receipientId, token.ProtocolVersion, token.SignedMessage)
	ecdsaV, err := subsig.NewECDSAVerifierFromPublicKey(GooglePaySHA256HashAlgorithm, GooglePayDEREncoding, publicKey)
	if err != nil {
		return err
	}
	return ecdsaV.Verify(signature, signedData)
}
