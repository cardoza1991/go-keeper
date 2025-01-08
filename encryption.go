package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"

    "golang.org/x/crypto/pbkdf2"
)

const (
    keyLength     = 32 // 256 bits
    saltLength    = 16
    iterations    = 100000
)

func deriveKey(masterPassword string, salt []byte) []byte {
    return pbkdf2.Key([]byte(masterPassword), salt, iterations, keyLength, sha256.New)
}

func encrypt(plaintext, masterPassword string) (string, error) {
    salt := make([]byte, saltLength)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(masterPassword, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
    combined := append(salt, ciphertext...)
    return hex.EncodeToString(combined), nil
}

func decrypt(encrypted, masterPassword string) (string, error) {
    combined, err := hex.DecodeString(encrypted)
    if err != nil {
        return "", err
    }

    if len(combined) < saltLength {
        return "", errors.New("invalid encrypted data")
    }

    salt := combined[:saltLength]
    ciphertext := combined[saltLength:]

    key := deriveKey(masterPassword, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("invalid ciphertext")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
