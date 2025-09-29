-- Crypto Module Examples
local crypto = require("crypto")

-- Hash functions
print("MD5:", crypto.md5("hello world"))
print("SHA256:", crypto.sha256("hello world"))
print("Generic hash:", crypto.hash("sha1", "hello world"))

-- Encoding/decoding
local encoded = crypto.base64_encode("hello world")
print("Base64 encoded:", encoded)
print("Base64 decoded:", crypto.base64_decode(encoded))

-- UUID generation
print("UUID:", crypto.uuid())
print("UUID v4:", crypto.uuid_v4())

-- Password generation
print("Password (8 chars):", crypto.generate_password(8, true))
print("Password (12 chars, no special):", crypto.generate_password(12, false))

-- AES encryption
local key = "my-secret-key-32-characters-long"
local plaintext = "sensitive data"
local encrypted = crypto.aes_encrypt(key, plaintext)
print("Encrypted:", encrypted)
print("Decrypted:", crypto.aes_decrypt(key, encrypted))

-- Random functions
print("Random string:", crypto.random_string(16))
print("Random bytes (base64):", crypto.random_bytes(16))