package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
)

// Credentials struct holds the session details obtained from the response
// of GetFederationToken call
type Credentials struct {
	AccessKeyId     string `json:"sessionId"`
	SecretAccessKey string `json:"sessionKey"`
	SessionToken    string `json:"sessionToken"`
}

// SigninTokenResp holds the response from token request call to AWS
type SigninTokenResp struct {
	SigninToken string `json:"SigninToken"`
}

func main() {
	user := flag.String("u", "test-user", "Name of temporary user")
	expiry := flag.Int64("x", 900, "Credentials expiration time in seconds")
	flag.Parse()

	// Create Session
	sess := session.Must(session.NewSession())

	// Create STS service
	svc := sts.New(sess, aws.NewConfig().WithRegion("us-west-2"))

        contents, err := ioutil.ReadFile("policy.json")
        if err != nil {
		log.Fatalf("Error reading policy file: %s", err)
	}

	// Create params for GetFederationToken()
	params := &sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(*expiry),
		Name:            aws.String(*user),
		Policy:          aws.String(string(contents)),
	}
	resp, err := svc.GetFederationToken(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	expiration := *resp.Credentials.Expiration

	// Create a Credentials struct from response and serialize
	sessionString, err := json.Marshal(Credentials{
		AccessKeyId:     *resp.Credentials.AccessKeyId,
		SecretAccessKey: *resp.Credentials.SecretAccessKey,
		SessionToken:    *resp.Credentials.SessionToken,
	})
	if err != nil {
		log.Fatalf("JSON marshal failed: %s", err)
	}

	// Build the federated token request URL
	tokenReqURL := buildTokenReqURL(string(sessionString))

	// Send the request to get an AWS console signin token
	signinToken := reqSigninToken(tokenReqURL.String())

	// Build the AWS console signin URL
	loginURL := buildLoginURL(signinToken)

	// Print credentials and console signin URL
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Printf("Access Key: ")
	fmt.Printf("%s\n", *resp.Credentials.AccessKeyId)
	c.Printf("Secret Key: ")
	fmt.Printf("%s\n", *resp.Credentials.SecretAccessKey)
	c.Printf("Session Token: ")
	fmt.Printf("%s\n", *resp.Credentials.SessionToken)
	c.Printf("URL: ")
	fmt.Printf("%s\n", loginURL)
	c.Printf("Credentials expires at: ")
	fmt.Printf("%s\n", expiration.Local().Format(time.UnixDate))
}

func reqSigninToken(reqURL string) string {
	var str SigninTokenResp
	res, err := http.Get(reqURL)
	if err != nil {
		log.Fatalf("HTTP request failed: %s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s", err)
	}
	err = json.Unmarshal(body, &str)
	if err != nil {
		log.Fatalf("JSON unmarshal failed: %s", err)
	}
	return str.SigninToken
}

func buildTokenReqURL(sessionString string) *url.URL {
	v := url.Values{}
	v.Add("Action", "getSigninToken")
	v.Add("Session", sessionString)
	return &url.URL{
		Scheme:   "https",
		Host:     "signin.aws.amazon.com",
		Path:     "federation",
		RawQuery: v.Encode(),
	}
}

func buildLoginURL(signinToken string) *url.URL {
	v := url.Values{}
	v.Add("Action", "login")
	v.Add("Destination", "https://console.aws.amazon.com/")
	v.Add("SigninToken", signinToken)
	return &url.URL{
		Scheme:   "https",
		Host:     "signin.aws.amazon.com",
		Path:     "federation",
		RawQuery: v.Encode(),
	}
}
