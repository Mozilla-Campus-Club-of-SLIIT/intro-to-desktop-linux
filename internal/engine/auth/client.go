package auth

var isAuthenticated bool = false

func VerifyAuth() bool {
	return isAuthenticated
}

func AuthUser() {
	isAuthenticated = true
}
