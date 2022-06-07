package signatures

import (
	"fmt"
	"reflect"

	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
)

// SignComponentDescriptor signs the given component-descriptor with the signer.
// The component-descriptor has to contain digests for componentReferences and resources.
func SignComponentDescriptor(cd *cdv2.ComponentDescriptor, signer Signer, hasher Hasher, signatureName string) error {
	hashedDigest, err := HashForComponentDescriptor(*cd, hasher)
	if err != nil {
		return fmt.Errorf("unable to get hash for component descriptor: %w", err)
	}

	signature, err := signer.Sign(*cd, *hashedDigest)
	if err != nil {
		return fmt.Errorf("unable to sign hash of normalised component descriptor: %w", err)
	}
	cd.Signatures = append(cd.Signatures, cdv2.Signature{
		Name:      signatureName,
		Digest:    *hashedDigest,
		Signature: *signature,
	})
	return nil
}

// VerifySignedComponentDescriptor verifies the signature (selected by signatureName) and hash of the component-descriptor (as specified in the signature).
// Does NOT resolve resources or referenced component-descriptors.
// Returns error if verification fails.
func VerifySignedComponentDescriptor(cd *cdv2.ComponentDescriptor, verifier Verifier, signatureName string) error {
	//find matching signature
	matchingSignature, err := GetSignatureByName(cd, signatureName)
	if err != nil {
		return fmt.Errorf("unable to get signature from component descriptor: %w", err)
	}

	//Verify author of signature
	err = verifier.Verify(*cd, *matchingSignature)
	if err != nil {
		return fmt.Errorf("unable to verify signature: %w", err)
	}

	//get hasher by algorithm name
	hasher, err := HasherForName(matchingSignature.Digest.HashAlgorithm)
	if err != nil {
		return fmt.Errorf("unable to create hasher for %s: %w", matchingSignature.Digest.HashAlgorithm, err)
	}

	//Verify normalised cd to given (and verified) hash
	calculatedDigest, err := HashForComponentDescriptor(*cd, *hasher)
	if err != nil {
		return fmt.Errorf("unable to hash component descriptor %s:%s: %w", cd.Name, cd.Version, err)
	}

	if !reflect.DeepEqual(*calculatedDigest, matchingSignature.Digest) {
		return fmt.Errorf("normalised component descriptor does not match hash from signature")
	}

	return nil
}

// GetSignatureByName returns the Signature (Digest and SigantureSpec) matching the given name
func GetSignatureByName(cd *cdv2.ComponentDescriptor, signatureName string) (*cdv2.Signature, error) {
	for _, signature := range cd.Signatures {
		if signature.Name == signatureName {
			return &signature, nil
		}
	}
	return nil, fmt.Errorf("signature with name %s not found in component descriptor", signatureName)

}
