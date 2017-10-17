package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	totp := flag.String("c", "", "OTP from MFA device")
	profile := flag.String("p", "default", "AWS credential profile")
	region := flag.String("r", "us-west-2", "AWS default region")
	expiry := flag.Int64("x", 86400, "Credentials expiration time in seconds")
	flag.Parse()

	// Create Session
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewSharedCredentials("", *profile),
	}))

	// Create STS service
	svc := sts.New(sess, aws.NewConfig().WithRegion(*region))

	output, err := svc.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: expiry,
		SerialNumber:    aws.String("arn:aws:iam::794111514033:mfa/slalom-simon.so"),
		TokenCode:       totp,
	})
	if err != nil {
		log.Fatalf("Error requesting temporary credentials: %s", err)
	}

	outputString := exportString("AWS_ACCESS_KEY_ID", *output.Credentials.AccessKeyId) +
		exportString("AWS_SECRET_ACCESS_KEY", *output.Credentials.SecretAccessKey) +
		exportString("AWS_SESSION_TOKEN", *output.Credentials.SessionToken) +
		exportString("AWS_DEFAULT_REGION", *region)
	fmt.Println(outputString)
}

func exportString(key string, val string) string {
	return fmt.Sprintf("export %s=%s;", key, val)
}
