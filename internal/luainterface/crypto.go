package luainterface

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
	"math/big"
	"strings"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

// CryptoModule provides cryptographic functionality for Lua scripts
type CryptoModule struct{}

// NewCryptoModule creates a new crypto module
func NewCryptoModule() *CryptoModule {
	return &CryptoModule{}
}

// RegisterCryptoModule registers the crypto module with the Lua state
func RegisterCryptoModule(L *lua.LState) {
	module := NewCryptoModule()
	
	// Create the crypto table
	cryptoTable := L.NewTable()
	
	// Hash functions
	L.SetField(cryptoTable, "md5", L.NewFunction(module.luaMD5))
	L.SetField(cryptoTable, "sha1", L.NewFunction(module.luaSHA1))
	L.SetField(cryptoTable, "sha256", L.NewFunction(module.luaSHA256))
	L.SetField(cryptoTable, "sha512", L.NewFunction(module.luaSHA512))
	L.SetField(cryptoTable, "hash", L.NewFunction(module.luaHash))
	
	// Encoding functions
	L.SetField(cryptoTable, "base64_encode", L.NewFunction(module.luaBase64Encode))
	L.SetField(cryptoTable, "base64_decode", L.NewFunction(module.luaBase64Decode))
	L.SetField(cryptoTable, "hex_encode", L.NewFunction(module.luaHexEncode))
	L.SetField(cryptoTable, "hex_decode", L.NewFunction(module.luaHexDecode))
	
	// UUID generation
	L.SetField(cryptoTable, "uuid", L.NewFunction(module.luaUUID))
	L.SetField(cryptoTable, "uuid_v4", L.NewFunction(module.luaUUIDv4))
	
	// Password generation
	L.SetField(cryptoTable, "generate_password", L.NewFunction(module.luaGeneratePassword))
	
	// AES encryption
	L.SetField(cryptoTable, "aes_encrypt", L.NewFunction(module.luaAESEncrypt))
	L.SetField(cryptoTable, "aes_decrypt", L.NewFunction(module.luaAESDecrypt))
	
	// Random functions
	L.SetField(cryptoTable, "random_string", L.NewFunction(module.luaRandomString))
	L.SetField(cryptoTable, "random_bytes", L.NewFunction(module.luaRandomBytes))
	
	// Register the crypto table globally
	L.SetGlobal("crypto", cryptoTable)
}

// Hash functions
func (c *CryptoModule) luaMD5(L *lua.LState) int {
	data := L.CheckString(1)
	hash := md5.Sum([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

func (c *CryptoModule) luaSHA1(L *lua.LState) int {
	data := L.CheckString(1)
	hash := sha1.Sum([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

func (c *CryptoModule) luaSHA256(L *lua.LState) int {
	data := L.CheckString(1)
	hash := sha256.Sum256([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

func (c *CryptoModule) luaSHA512(L *lua.LState) int {
	data := L.CheckString(1)
	hash := sha512.Sum512([]byte(data))
	L.Push(lua.LString(hex.EncodeToString(hash[:])))
	return 1
}

func (c *CryptoModule) luaHash(L *lua.LState) int {
	algorithm := L.CheckString(1)
	data := L.CheckString(2)
	
	var hash []byte
	switch strings.ToLower(algorithm) {
	case "md5":
		h := md5.Sum([]byte(data))
		hash = h[:]
	case "sha1":
		h := sha1.Sum([]byte(data))
		hash = h[:]
	case "sha256":
		h := sha256.Sum256([]byte(data))
		hash = h[:]
	case "sha512":
		h := sha512.Sum512([]byte(data))
		hash = h[:]
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("unsupported hash algorithm: " + algorithm))
		return 2
	}
	
	L.Push(lua.LString(hex.EncodeToString(hash)))
	return 1
}

// Encoding functions
func (c *CryptoModule) luaBase64Encode(L *lua.LState) int {
	data := L.CheckString(1)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	L.Push(lua.LString(encoded))
	return 1
}

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

func (c *CryptoModule) luaHexEncode(L *lua.LState) int {
	data := L.CheckString(1)
	encoded := hex.EncodeToString([]byte(data))
	L.Push(lua.LString(encoded))
	return 1
}

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

// UUID functions
func (c *CryptoModule) luaUUID(L *lua.LState) int {
	id := uuid.New()
	L.Push(lua.LString(id.String()))
	return 1
}

func (c *CryptoModule) luaUUIDv4(L *lua.LState) int {
	id := uuid.New()
	L.Push(lua.LString(id.String()))
	return 1
}

// Password generation
func (c *CryptoModule) luaGeneratePassword(L *lua.LState) int {
	length := L.CheckInt(1)
	includeSpecial := L.OptBool(2, true)
	
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)
	
	charset := lowercase + uppercase + digits
	if includeSpecial {
		charset += special
	}
	
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		password[i] = charset[n.Int64()]
	}
	
	L.Push(lua.LString(string(password)))
	return 1
}

// AES encryption/decryption
func (c *CryptoModule) luaAESEncrypt(L *lua.LState) int {
	key := L.CheckString(1)
	plaintext := L.CheckString(2)
	
	// Ensure key is 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		hash := sha256.Sum256(keyBytes)
		keyBytes = hash[:]
	}
	
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	// Generate random IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	L.Push(lua.LString(encoded))
	return 1
}

func (c *CryptoModule) luaAESDecrypt(L *lua.LState) int {
	key := L.CheckString(1)
	ciphertextB64 := L.CheckString(2)
	
	// Ensure key is 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		hash := sha256.Sum256(keyBytes)
		keyBytes = hash[:]
	}
	
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	if len(ciphertext) < aes.BlockSize {
		L.Push(lua.LNil)
		L.Push(lua.LString("ciphertext too short"))
		return 2
	}
	
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	
	L.Push(lua.LString(string(ciphertext)))
	return 1
}

// Random functions
func (c *CryptoModule) luaRandomString(L *lua.LState) int {
	length := L.CheckInt(1)
	charset := L.OptString(2, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		result[i] = charset[n.Int64()]
	}
	
	L.Push(lua.LString(string(result)))
	return 1
}

func (c *CryptoModule) luaRandomBytes(L *lua.LState) int {
	length := L.CheckInt(1)
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	encoded := base64.StdEncoding.EncodeToString(bytes)
	L.Push(lua.LString(encoded))
	return 1
}