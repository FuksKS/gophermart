package crypto

type PasswordEncrypter interface {
	PassEncrypt(plaintext string) (string, error)
	PassDecrypt(encrypted string) (string, error)
}
