package archive

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize  = 16
	nonceSize = 12
)

// An EncryptedArchive represents an encrypted archive.
type EncryptedArchive struct {
	bytes []byte
	salt  []byte
	nonce []byte
}

// Write writes the encrypted archive into the provided writer.
// The encrypted archive is serialized as follows:
//
//  1. The magic field is serialized as a 4-byte sequence.
//  2. The salt field is serialized as a 16-byte sequence.
//  3. The nonce field is serialized as a sequence of bytes.
//  4. The encrypted data is serialized as a sequence of bytes.
func (a *EncryptedArchive) Write(w io.Writer) error {
	// Write the magic (4 bytes)
	if _, err := w.Write(encMagic); err != nil {
		return err
	}

	// Write the salt (16 bytes)
	if _, err := w.Write(a.salt); err != nil {
		return err
	}

	// Write the nonce
	if _, err := w.Write(a.nonce); err != nil {
		return err
	}

	// Write the encrypted data
	if _, err := w.Write(a.bytes); err != nil {
		return err
	}

	return nil
}

// ReadEncryptedArchive reads an encrypted archive from the provided reader.
func ReadEncryptedArchive(r io.Reader) (*EncryptedArchive, error) {
	if err := mustReadEncryptedMagic(r); err != nil {
		return nil, err
	}

	// Read the salt (16 bytes)
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(r, salt); err != nil {
		return nil, err
	}

	// Read the nonce
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	// Read the encrypted data
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &EncryptedArchive{
		bytes: data,
		salt:  salt,
		nonce: nonce,
	}, nil
}

// Encrypt encrypts the archive using AES-GCM with the provided password.
func (a *Archive) Encrypt(password string) (*EncryptedArchive, error) {
	// Generate a salt for key derivation (PBKDF2)
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	aesGCM, err := newCipher(password, salt)
	if err != nil {
		return nil, err
	}

	// Generate a nonce for AES-GCM (random IV)
	// aesGCM.NonceSize() returns 12 bytes
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Get the plaintext data
	plaintext, err := a.GetBytes()
	if err != nil {
		return nil, err
	}

	// Encrypt the data using AES-GCM
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	return &EncryptedArchive{
		bytes: ciphertext,
		salt:  salt,
		nonce: nonce,
	}, nil
}

// Decrypt decrypts the encrypted archive using AES-GCM with the provided password.
// If the password is incorrect, the process will fail as the Archive data will be
// corrupted.
func (a *EncryptedArchive) Decrypt(password string) (*Archive, error) {
	aesGCM, err := newCipher(password, a.salt)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, a.nonce, a.bytes, nil)
	if err != nil {
		return nil, err
	}

	return ReadArchive(bytes.NewReader(plaintext))
}

// newCipher creates a new AES-GCM cipher with the provided password and salt.
func newCipher(password string, salt []byte) (cipher.AEAD, error) {
	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}
