package tls

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BoundSASigningKey contains a user provided key and public parts for the
// service account signing key used by kube-apiserver.
// This asset does not generate any new content and only loads these files from disk
// when provided by the user.
type BoundSASigningKey struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*BoundSASigningKey)(nil)

// Name returns a human friendly name for the asset.
func (*BoundSASigningKey) Name() string {
	return "User-provided Service Account Signing key"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*BoundSASigningKey) Dependencies() []asset.Asset {
	return nil
}

// Generate generates the CloudProviderConfig.
func (*BoundSASigningKey) Generate(dependencies asset.Parents) error { return nil }

// Files returns the files generated by the asset.
func (sk *BoundSASigningKey) Files() []*asset.File {
	return sk.FileList
}

// Load reads the private key from the disk.
// It ensures that the key provided is a valid RSA key.
func (sk *BoundSASigningKey) Load(f asset.FileFetcher) (bool, error) {
	keyFile, err := f.FetchByName(filepath.Join(tlsDir, "bound-service-account-signing-key.key"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	rsaKey, err := PemToPrivateKey(keyFile.Data)
	if err != nil {
		logrus.Debugf("Failed to load rsa.PrivateKey from file: %s", err)
		return false, errors.Wrap(err, "failed to load rsa.PrivateKey from the file")
	}
	pubData, err := PublicKeyToPem(&rsaKey.PublicKey)
	if err != nil {
		return false, errors.Wrap(err, "failed to extract public key from the key")
	}
	sk.FileList = []*asset.File{keyFile, {Filename: filepath.Join(tlsDir, "bound-service-account-signing-key.pub"), Data: pubData}}
	return true, nil
}