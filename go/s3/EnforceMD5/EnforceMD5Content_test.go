// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/google/uuid"
)

type Config struct {
	Bucket string `json:"Bucket"`
	Key    string `json:"Key"`
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration(t *testing.T) error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	t.Log("Bucket: " + globalConfig.Bucket)
	t.Log("Key:    " + globalConfig.Key)

	return nil
}

func createBucket(sess *session.Session, bucket *string) error {
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteBucket(sess *session.Session, bucket *string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		return err
	}

	return nil
}

func TestEnforceMD5Content(t *testing.T) {
	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	err := populateConfiguration(t)
	if err != nil {
		t.Fatal(err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	bucketCreated := false

	if globalConfig.Bucket == "" || globalConfig.Key == "" {
		id := uuid.New()
		globalConfig.Bucket = "test-bucket-" + id.String()
		globalConfig.Key = "test-key.txt"

		err := createBucket(sess, &globalConfig.Bucket)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Created bucket " + globalConfig.Bucket)
	}

	defClient, err := SetMd5(sess, &globalConfig.Bucket, &globalConfig.Key)
	if err != nil {
		t.Log("You'll have to delete " + globalConfig.Bucket + " yourself")
		t.Fatal(err)
	}

	t.Log(defClient, err)

	if bucketCreated {
		err := deleteBucket(sess, &globalConfig.Bucket)
		if err != nil {
			t.Log("You'll have to delete " + globalConfig.Bucket + " yourself")
			t.Fatal(err)
		}

		t.Log("Deleted " + globalConfig.Bucket)
	}
}
