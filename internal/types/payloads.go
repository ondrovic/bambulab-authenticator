package types

type LoginPayload struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	ApiError string `json:"apiError"`
}

type TwoFactorPayload struct {
	TFAKey  string `json:"tfaKey"`
	TFACode string `json:"tfaCode"`
}

type RequestEmailCodePayload struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type EmailCodePayload struct {
	Account string `json:"account"`
	Code    string `json:"code"`
}

type LoginResponse struct {
	AccessToken      string `json:"accessToken,omitempty"`
	RefreshToken     string `json:"refreshToken,omitempty"`
	ExpiresIn        int    `json:"expiresIn,omitempty"`
	RefreshExpiresIn int    `json:"refreshExpiresIn,omitempty"`
	TfaKey           string `json:"tfaKey,omitempty"`
	AccessMethod     string `json:"accessMethod,omitempty"`
	LoginType        string `json:"loginType,omitempty"`
}
