package auth

import (
	"encoding/json"
	"fmt"

	"github.com/ondrovic/bambulab-authenticator/internal/consts"
	"github.com/ondrovic/bambulab-authenticator/internal/httpclient"
	"github.com/ondrovic/bambulab-authenticator/internal/types"
	"github.com/ondrovic/bambulab-authenticator/internal/utils"
)

func Login(opts *types.CliFlags) error {

	if httpclient.Client == nil {
		if err := httpclient.InitClient(consts.EMPTY_STRING); err != nil {
			return err
		}
	}

	loginPayload := types.LoginPayload{
		Account:  opts.UserAccount,
		Password: opts.UserPassword,
		ApiError: consts.EMPTY_STRING,
	}

	jsonLoginPayload, err := json.Marshal(loginPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal loginPayload: %v", err)
	}

	url, err := consts.RegionalURL(consts.LoginURL, opts.UserRegion)
	if err != nil {
		return fmt.Errorf("failed to construct regionalUrl: %v", err)
	}

	resp, err := httpclient.Request("POST", string(url), jsonLoginPayload)
	if err != nil {
		return err
	}

	processLoginType(resp, opts)

	return nil
}

func processLoginType(loginResponse *types.LoginResponse, opts *types.CliFlags) error {
	switch loginResponse.LoginType {
	case "verifyCode":
		if err := sendCodeToEmail(opts); err != nil {
			fmt.Printf("error sending email %v\n", err)
			return err
		}

		fmt.Print("VerifyCode: Enter the code from your email: ")
		var verifyCode string
		fmt.Scanln(&verifyCode)

		if err := emailCodeLogin(verifyCode, opts); err != nil {
			return err
		}

		return nil
	case "tfa":
		twoFactorAuth(loginResponse.TfaKey, opts)
		return nil
	default:
		return fmt.Errorf("unknown login type: %v", loginResponse.LoginType)
	}
}

func sendCodeToEmail(opts *types.CliFlags) error {

	sendCodePayload := types.RequestEmailCodePayload{
		Email: opts.UserAccount,
		Type:  "codeLogin",
	}

	jsonSendCodePayload, err := json.Marshal(sendCodePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal sendCodePayload: %v", err)
	}

	url, err := consts.RegionalURL(consts.EmailCodeURL, opts.UserRegion)
	if err != nil {
		return fmt.Errorf("failed to construct emailCodeUrl: %v", err)
	}

	_, err = httpclient.Request("POST", string(url), jsonSendCodePayload)

	if err != nil {
		return err
	}

	return nil
}

func emailCodeLogin(code string, opts *types.CliFlags) error {

	emailCodePayload := types.EmailCodePayload{
		Account: opts.UserAccount,
		Code:    code,
	}

	jsonEmailCodePayload, err := json.Marshal(emailCodePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal emailCodePayload: %v", err)
	}

	url, err := consts.RegionalURL(consts.LoginURL, opts.UserRegion)
	if err != nil {
		return fmt.Errorf("failed to construct twoFactorUrl: %v", err)
	}

	emailCodeResponse, err := httpclient.Request("POST", string(url), jsonEmailCodePayload)
	if err != nil {
		return err
	}

	if err := utils.SaveLoginResponseToFile(*emailCodeResponse, opts.OutputPath); err != nil {
		return err
	}

	return nil
}

func twoFactorAuth(tfaKey string, opts *types.CliFlags) error {
	fmt.Print("2FA: Enter your one-time password: ")

	var tfaCode string
	fmt.Scanln(&tfaCode)

	twoFactorAuthPayload := types.TwoFactorPayload{
		TFAKey:  tfaKey,
		TFACode: tfaCode,
	}

	twoFactorAuthPayloadJSON, err := json.Marshal(twoFactorAuthPayload)
	if err != nil {
		return fmt.Errorf("unable to marshal twoFactorAuthPayload: %v", err)
	}

	url, err := consts.RegionalURL(consts.TwoFactorURL, opts.UserRegion)
	if err != nil {
		return fmt.Errorf("failed to construct twoFactorUrl: %v", err)
	}

	tfaResponse, err := httpclient.CookieRequest("POST", string(url), twoFactorAuthPayloadJSON)

	if err := utils.SaveLoginResponseToFile(*tfaResponse, opts.OutputPath); err != nil {
		return err
	}

	return nil
}
