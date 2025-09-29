package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/yuin/gopher-lua"
	"golang.org/x/crypto/bcrypt"
)

// CryptoModule provides cryptographic operations
type CryptoModule struct {
	info CoreModuleInfo
}

// NewCryptoModule creates a new crypto module
func NewCryptoModule() *CryptoModule {
	info := CoreModuleInfo{
		Name:         "crypto",
		Version:      "1.0.0",
		Description:  "Cryptographic operations including hashing, encryption, and encoding",
		Author:       "Sloth Runner Team",
		Category:     "core",
		Dependencies: []string{},
	}

	return &CryptoModule{
		info: info,
	}
}

// Info returns module information
func (c *CryptoModule) Info() CoreModuleInfo {
	return c.info
}

// Loader loads the crypto module into Lua
func (c *CryptoModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"md5":           c.luaMD5,
		"sha1":          c.luaSHA1,
		"sha256":        c.luaSHA256,
		"sha512":        c.luaSHA512,
		"bcrypt_hash":   c.luaBcryptHash,
		"bcrypt_check":  c.luaBcryptCheck,
		"base64_encode": c.luaBase64Encode,
		"base64_decode": c.luaBase64Decode,
		"hex_encode":    c.luaHexEncode,
		"hex_decode":    c.luaHexDecode,
		"aes_encrypt":   c.luaAESEncrypt,
		"aes_decrypt":   c.luaAESDecrypt,
		"random_bytes":  c.luaRandomBytes,
		"random_string": c.luaRandomString,
	})

	L.Push(mod)
	return 1
}

// luaMD5 generates MD5 hash
func (c *CryptoModule) luaMD5(L *lua.LState) int {
	data := L.CheckString(1)
	
	hash := md5.Sum([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

// luaSHA1 generates SHA1 hash
func (c *CryptoModule) luaSHA1(L *lua.LState) int {
	data := L.CheckString(1)
	
	hash := sha1.Sum([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

// luaSHA256 generates SHA256 hash
func (c *CryptoModule) luaSHA256(L *lua.LState) int {
	data := L.CheckString(1)
	
	hash := sha256.Sum256([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

// luaSHA512 generates SHA512 hash
func (c *CryptoModule) luaSHA512(L *lua.LState) int {
	data := L.CheckString(1)
	
	hash := sha512.Sum512([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

// luaBcryptHash generates bcrypt hash
func (c *CryptoModule) luaBcryptHash(L *lua.LState) int {
	password := L.CheckString(1)
	cost := 10 // default cost
	
	if L.GetTop() > 1 {
		cost = int(L.CheckNumber(2))
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(hash)))
	return 1
}

// luaBcryptCheck verifies bcrypt hash
func (c *CryptoModule) luaBcryptCheck(L *lua.LState) int {
	password := L.CheckString(1)
	hash := L.CheckString(2)
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	L.Push(lua.LBool(err == nil))
	return 1
}

// luaBase64Encode encodes data to base64
func (c *CryptoModule) luaBase64Encode(L *lua.LState) int {
	data := L.CheckString(1)
	
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	L.Push(lua.LString(encoded))
	return 1
}

// luaBase64Decode decodes base64 data
func (c *CryptoModule) luaBase64Decode(L *lua.LState) int {
	data := L.CheckString(1)
	
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(decoded)))
	return 1
}

// luaHexEncode encodes data to hexadecimal
func (c *CryptoModule) luaHexEncode(L *lua.LState) int {
	data := L.CheckString(1)
	
	encoded := hex.EncodeToString([]byte(data))
	L.Push(lua.LString(encoded))
	return 1
}

// luaHexDecode decodes hexadecimal data
func (c *CryptoModule) luaHexDecode(L *lua.LState) int {
	data := L.CheckString(1)
	
	decoded, err := hex.DecodeString(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(decoded)))
	return 1
}

// luaAESEncrypt encrypts data using AES
func (c *CryptoModule) luaAESEncrypt(L *lua.LState) int {
	plaintext := L.CheckString(1)
	key := L.CheckString(2)
	
	// Ensure key is 32 bytes for AES-256
	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))
	
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// Encode to base64 for easy handling
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	L.Push(lua.LString(encoded))
	return 1
}

// luaAESDecrypt decrypts AES encrypted data
func (c *CryptoModule) luaAESDecrypt(L *lua.LState) int {
	ciphertext := L.CheckString(1)
	key := L.CheckString(2)
	
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Ensure key is 32 bytes for AES-256
	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))
	
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		L.Push(lua.LNil)
		L.Push(lua.LString("ciphertext too short"))
		return 2
	}
	
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(plaintext)))
	return 1
}

// luaRandomBytes generates random bytes
func (c *CryptoModule) luaRandomBytes(L *lua.LState) int {
	size := int(L.CheckNumber(1))
	
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(hex.EncodeToString(bytes)))
	return 1
}

// luaRandomString generates a random string
func (c *CryptoModule) luaRandomString(L *lua.LState) int {
	length := int(L.CheckNumber(1))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	if L.GetTop() > 1 {
		charset = L.CheckString(2)
	}
	
	bytes := make([]byte, length)
	for i := range bytes {
		randomIndex := make([]byte, 1)
		rand.Read(randomIndex)
		bytes[i] = charset[int(randomIndex[0])%len(charset)]
	}
	
	L.Push(lua.LString(string(bytes)))
	return 1
}