package core

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestCryptoModule_Info(t *testing.T) {
	module := NewCryptoModule()
	info := module.Info()

	if info.Name != "crypto" {
		t.Errorf("Expected module name 'crypto', got '%s'", info.Name)
	}

	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}
}

func TestCryptoModule_MD5(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local hash = crypto.md5("hello world")
		return hash
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"
	if result.String() != expected {
		t.Errorf("Expected MD5 hash '%s', got '%s'", expected, result.String())
	}
}

func TestCryptoModule_SHA1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local hash = crypto.sha1("hello world")
		return hash
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	expected := "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
	if result.String() != expected {
		t.Errorf("Expected SHA1 hash '%s', got '%s'", expected, result.String())
	}
}

func TestCryptoModule_SHA256(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local hash = crypto.sha256("hello world")
		return hash
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if result.String() != expected {
		t.Errorf("Expected SHA256 hash '%s', got '%s'", expected, result.String())
	}
}

func TestCryptoModule_SHA512(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local hash = crypto.sha512("hello")
		return hash
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	// Just check it's a valid SHA512 hash (128 hex chars)
	if len(result.String()) != 128 {
		t.Errorf("Expected SHA512 hash length 128, got %d", len(result.String()))
	}
}

func TestCryptoModule_BcryptHashAndCheck(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local password = "my_secret_password"
		local hash = crypto.bcrypt_hash(password, 4)
		
		if not hash then
			error("Failed to generate bcrypt hash")
		end
		
		-- Verify correct password
		local valid = crypto.bcrypt_check(password, hash)
		if not valid then
			error("Password verification failed")
		end
		
		-- Verify wrong password
		local invalid = crypto.bcrypt_check("wrong_password", hash)
		if invalid then
			error("Wrong password should not match")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestCryptoModule_Base64(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local original = "Hello, World!"
		local encoded = crypto.base64_encode(original)
		local decoded = crypto.base64_decode(encoded)
		
		if decoded ~= original then
			error("Base64 encode/decode failed: " .. decoded .. " ~= " .. original)
		end
		
		return encoded
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	expected := "SGVsbG8sIFdvcmxkIQ=="
	if result.String() != expected {
		t.Errorf("Expected base64 '%s', got '%s'", expected, result.String())
	}
}

func TestCryptoModule_Hex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local original = "test123"
		local encoded = crypto.hex_encode(original)
		local decoded = crypto.hex_decode(encoded)
		
		if decoded ~= original then
			error("Hex encode/decode failed")
		end
		
		return encoded
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	expected := "74657374313233"
	if result.String() != expected {
		t.Errorf("Expected hex '%s', got '%s'", expected, result.String())
	}
}

func TestCryptoModule_AESEncryptDecrypt(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local plaintext = "This is a secret message!"
		local key = "my-secret-encryption-key-32bytes"
		
		local encrypted = crypto.aes_encrypt(plaintext, key)
		if not encrypted then
			error("Encryption failed")
		end
		
		local decrypted = crypto.aes_decrypt(encrypted, key)
		if not decrypted then
			error("Decryption failed")
		end
		
		if decrypted ~= plaintext then
			error("Decrypted text doesn't match original: " .. decrypted)
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestCryptoModule_RandomBytes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local random1 = crypto.random_bytes(16)
		local random2 = crypto.random_bytes(16)
		
		if #random1 ~= 32 then  -- 16 bytes = 32 hex chars
			error("Expected 32 hex characters, got " .. #random1)
		end
		
		if random1 == random2 then
			error("Random bytes should be different")
		end
		
		return random1
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	if len(result.String()) != 32 {
		t.Errorf("Expected 32 hex characters, got %d", len(result.String()))
	}
}

func TestCryptoModule_RandomString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local random1 = crypto.random_string(20)
		local random2 = crypto.random_string(20)
		
		if #random1 ~= 20 then
			error("Expected length 20, got " .. #random1)
		end
		
		if random1 == random2 then
			error("Random strings should be different")
		end
		
		-- Test custom charset
		local numbers = crypto.random_string(10, "0123456789")
		if #numbers ~= 10 then
			error("Expected length 10, got " .. #numbers)
		end
		
		return random1
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}

	result := L.Get(-1)
	if len(result.String()) != 20 {
		t.Errorf("Expected length 20, got %d", len(result.String()))
	}
}

func TestCryptoModule_InvalidBase64Decode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local decoded, err = crypto.base64_decode("invalid!!base64!!")
		
		if decoded ~= nil then
			error("Expected nil for invalid base64")
		end
		
		if not err then
			error("Expected error message")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}

func TestCryptoModule_InvalidHexDecode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	module := NewCryptoModule()
	L.PreloadModule("crypto", module.Loader)

	code := `
		local crypto = require("crypto")
		local decoded, err = crypto.hex_decode("invalid hex gg")
		
		if decoded ~= nil then
			error("Expected nil for invalid hex")
		end
		
		if not err then
			error("Expected error message")
		end
		
		return "success"
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Lua execution failed: %v", err)
	}
}
