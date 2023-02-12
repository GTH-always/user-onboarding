package helpers

import (
	"context"
	"fmt"
	"user-onboarding/config"
	dynamo "user-onboarding/services/s3Bucket"
	structs "user-onboarding/struct"
	requestStruct "user-onboarding/struct/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/getsentry/sentry-go"
	"golang.org/x/crypto/bcrypt"
)

func UserDetails(ctx context.Context, request *structs.UserDetails, sentryCtx context.Context) error {
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[DAO] UserDetails") //sentry to log db calls
	defer span.Finish()
	svc := dynamodb.New(dynamo.AwsSession())

	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] Check if user data is present")

	email := structs.UserDetails{
		Email: request.Email,
	}
	key, err := dynamodbattribute.MarshalMap(email)
	if err != nil {
		return err
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(config.Get().Table),
	}
	result, err := svc.GetItem(input)
	fmt.Println(result)
	dbSpan1.Finish() //noting time of query
	if err != nil || result.Item != nil {
		err = fmt.Errorf("email already exists")
		return err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 1)
	request.Password = string(hashedPassword)

	key, err = dynamodbattribute.MarshalMap(request)

	if err != nil {
		return err
	}
	dbSpan2 := sentry.StartSpan(span.Context(), "[DB] Inserting user data in table")

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(config.Get().Table),
		Item:      key,
	})

	dbSpan2.Finish()

	if err != nil {
		return err
	}

	return nil
}

func FetchUser(ctx context.Context, request *structs.UserDetails, sentryCtx context.Context) (map[string]*dynamodb.AttributeValue, error) {
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[DAO] UserDetails") //sentry to log db calls
	defer span.Finish()
	svc := dynamodb.New(dynamo.AwsSession())

	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] Check if user data is present")
	email := structs.UserDetails{
		Email: request.Email,
	}

	key, err := dynamodbattribute.MarshalMap(email)

	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(config.Get().Table),
	}
	result, err := svc.GetItem(input)

	dbSpan1.Finish() //noting time of query
	if err != nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return result.Item, err
}

func UserLogin(ctx context.Context, request *requestStruct.UserLogin, sentryCtx context.Context) error {
	defer sentry.Recover()
	span := sentry.StartSpan(sentryCtx, "[DAO] Userlogin") //sentry to log db calls
	defer span.Finish()
	svc := dynamodb.New(dynamo.AwsSession())

	dbSpan1 := sentry.StartSpan(span.Context(), "[DB] User login")
	email := structs.UserDetails{
		Email: request.Email,
	}

	key, err := dynamodbattribute.MarshalMap(email)

	if err != nil {
		return err
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(config.Get().Table),
	}
	result, err := svc.GetItem(input)

	if err != nil {
		fmt.Println(err)
		return err
	}

	userPassword := ""

	for key, v := range result.Item {
		if key == "password" {
			userPassword = *v.S
			break
		}
	}
	dbSpan1.Finish() //noting time of query
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(request.Password))

	if err != nil {
		return err
	}
	return err
}
