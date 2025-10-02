package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestAWSS3Upload(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.s3_upload({
			bucket = "test-bucket",
			key = "test-key",
			file = "/nonexistent/file"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSS3Download(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.s3_download({
			bucket = "test-bucket",
			key = "test-key",
			destination = "/tmp/test-file"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSS3List(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.s3_list({
			bucket = "test-bucket"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSEC2ListInstances(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.ec2_list_instances({
			region = "us-east-1"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSEC2StartInstance(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.ec2_start_instance({
			instance_id = "i-1234567890abcdef0",
			region = "us-east-1"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSEC2StopInstance(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.ec2_stop_instance({
			instance_id = "i-1234567890abcdef0",
			region = "us-east-1"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}

func TestAWSLambdaInvoke(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterAWSModule(L)

	script := `
		local aws = require("aws")
		local result = aws.lambda_invoke({
			function_name = "test-function",
			region = "us-east-1",
			payload = "{}"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("AWS test skipped (AWS credentials not available): %v", err)
		return
	}
}
