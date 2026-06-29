package userpool

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoIDPTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type Client interface {
	AddUser(ctx context.Context, email, password string) error
	DeleteUser(ctx context.Context, email string) error
}

func NewClient(cognitoClient *cognitoidentityprovider.Client, userPoolId string) Client {
	return &client{cognitoClient: cognitoClient, userPoolId: userPoolId}
}

type client struct {
	cognitoClient *cognitoidentityprovider.Client
	userPoolId    string
}

func (c *client) AddUser(ctx context.Context, email, password string) error {
	_, err := c.cognitoClient.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: &c.userPoolId,
		Username:   &email,
		UserAttributes: []cognitoIDPTypes.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
		DesiredDeliveryMediums: []cognitoIDPTypes.DeliveryMediumType{},
		MessageAction:          cognitoIDPTypes.MessageActionTypeSuppress,
	})
	if err != nil {
		return err
	}

	_, err = c.cognitoClient.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: &c.userPoolId,
		Username:   &email,
		Password:   &password,
		Permanent:  true,
	})
	return err
}

func (c *client) DeleteUser(ctx context.Context, email string) error {
	_, err := c.cognitoClient.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: &c.userPoolId,
		Username:   &email,
	})
	if _, ok := errors.AsType[*cognitoIDPTypes.UserNotFoundException](err); ok {
		return nil
	}
	return err
}

func RandomPassword(length int) string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!@#$%^&*()-_=+[]{}|:,.?/"
	all := lower + upper + digits + special

	randInt := func(max int) int {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			return 0
		}
		return int(nBig.Int64())
	}

	buf := make([]byte, 0, length)
	buf = append(buf, lower[randInt(len(lower))])
	buf = append(buf, upper[randInt(len(upper))])
	buf = append(buf, digits[randInt(len(digits))])
	buf = append(buf, special[randInt(len(special))])

	for len(buf) < length {
		buf = append(buf, all[randInt(len(all))])
	}

	for i := len(buf) - 1; i > 0; i-- {
		j := randInt(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}
