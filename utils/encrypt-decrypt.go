package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)


const saltSize = 16


//GenerateRandomSalt ...
func GenerateRandomSalt(saltSize int) []byte {
  var salt = make([]byte, saltSize)

  _, err := rand.Read(salt[:])

  if err != nil {
    panic(err)
  }

  return salt
}

//HashPassword ...
func HashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha256Hasher = sha256.New()
	passwordBytes = append(passwordBytes, salt...)
	sha256Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha256Hasher.Sum(nil)
	var base64EncodedPasswordHash =
		base64.URLEncoding.EncodeToString(hashedPasswordBytes)
  
	return base64EncodedPasswordHash
  }


  //DoPasswordsMatch ...
  func DoPasswordsMatch(hashedPassword, currPassword string, salt[]byte) bool {
	var currPasswordHash = HashPassword(currPassword, salt)
	return hashedPassword == currPasswordHash
  }
  