package consts

import (
	"errors"
	"strings"

	"github.com/ondrovic/bambulab-authenticator/internal/utils"
)

type URL string

const (
	EmailCodeURL URL = "https://api.bambulab.com/v1/user-service/user/sendemail/code"
	LoginURL     URL = "https://api.bambulab.com/v1/user-service/user/login"
	ProfileURL   URL = "https://api.bambulab.com/v1/user-service/my/profile"
	RefererURL   URL = "https://bambulab.com"
	TwoFactorURL URL = "https://bambulab.com/api/sign-in/tfa"
)

func RegionalURL(url URL, region string) (URL, error) {

	if utils.IsEmpty(region) {
		return "", errors.New("region cannot be empty")
	}

	region = strings.ToLower(region)

	if region == "china" {
		regionalUrl := strings.Replace(string(url), ".com", ".cn", -1)
		return URL(regionalUrl), nil
	}

	return url, nil
}
