package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MythicMeta/MythicContainer"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
)

const VERSION = "2026.6.1"

var compressorDefinition = agentstructs.PayloadType{
	Name:                   "compressor",
	Description:            "Compress a payload using a variety of methods",
	FileExtension:          "zip",
	Wrapper:                true,
	SupportsDynamicLoading: false,
	Author:                 "@lum8rjack",
	SemVer:                 VERSION,
	SupportedOS:            []string{agentstructs.SUPPORTED_OS_WINDOWS, agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_MACOS},
	AgentType:              agentstructs.AgentTypeWrapper,
	BuildParameters: []agentstructs.BuildParameter{
		{
			Name:          "name",
			Description:   "Specify the file name of the payload inside the compressed archive",
			Required:      true,
			DefaultValue:  "agent.bin",
			ParameterType: agentstructs.BUILD_PARAMETER_TYPE_STRING,
			UiPosition:    1,
			GroupName:     "Method",
			SupportedOS:   []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_WINDOWS, agentstructs.SUPPORTED_OS_MACOS},
		},
		{
			Name:          "method",
			Description:   "Specify the method to use to compress the payload",
			Required:      true,
			DefaultValue:  "zip",
			Choices:       []string{"zip", "tar", "tar.gz", "tar.bz2", "tar.xz"},
			ParameterType: agentstructs.BUILD_PARAMETER_TYPE_CHOOSE_ONE,
			UiPosition:    2,
			GroupName:     "Method",
			SupportedOS:   []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_WINDOWS, agentstructs.SUPPORTED_OS_MACOS},
		},
		{
			Name:          "use_password",
			Description:   "Use a password to encrypt the payload",
			Required:      false,
			DefaultValue:  false,
			ParameterType: agentstructs.BUILD_PARAMETER_TYPE_BOOLEAN,
			UiPosition:    3,
			GroupName:     "Password",
			SupportedOS:   []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_WINDOWS, agentstructs.SUPPORTED_OS_MACOS},
			HideConditions: []agentstructs.BuildParameterHideCondition{
				{Name: "method", Operand: agentstructs.HideConditionOperandNotEQ, Value: "zip"},
			},
		},
		{
			Name:          "password",
			Description:   "The password to encrypt the payload",
			Required:      true,
			DefaultValue:  "",
			ParameterType: agentstructs.BUILD_PARAMETER_TYPE_STRING,
			UiPosition:    4,
			GroupName:     "Password",
			SupportedOS:   []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_WINDOWS, agentstructs.SUPPORTED_OS_MACOS},
			HideConditions: []agentstructs.BuildParameterHideCondition{
				{Name: "use_password", Operand: agentstructs.HideConditionOperandEQ, Value: "false"},
			},
		},
	},
}

func build(payloadBuildMsg agentstructs.PayloadBuildMessage) agentstructs.PayloadBuildResponse {
	payloadBuildResponse := agentstructs.PayloadBuildResponse{
		PayloadUUID: *payloadBuildMsg.WrappedPayloadUUID,
		Success:     true,
	}

	// Setup logger for the docker container
	prefix := fmt.Sprintf("[builder:%s]", *payloadBuildMsg.WrappedPayloadUUID)
	customLogger := log.New(os.Stdout, prefix, log.Default().Flags())

	outputName := "package"

	// Get build arguments
	payloadName, err := payloadBuildMsg.BuildParameters.GetStringArg("name")
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = err.Error()
		customLogger.Println("Failed to get build argument: name")
		return payloadBuildResponse
	}

	method, err := payloadBuildMsg.BuildParameters.GetStringArg("method")
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = err.Error()
		customLogger.Println("Failed to get build argument: method")
		return payloadBuildResponse
	}

	usePassword, err := payloadBuildMsg.BuildParameters.GetBooleanArg("use_password")
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = err.Error()
		customLogger.Println("Failed to get build argument: use_password")
		return payloadBuildResponse
	}

	password, err := payloadBuildMsg.BuildParameters.GetStringArg("password")
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = err.Error()
		customLogger.Println("Failed to get build argument: password")
		return payloadBuildResponse
	}

	// Check method
	command := ""
	finalOutputName := ""

	if method == "zip" {
		finalOutputName = fmt.Sprintf("%s.zip", outputName)
		command = fmt.Sprintf("zip -j %s %s", finalOutputName, payloadName)
		if usePassword {
			command = fmt.Sprintf("zip -jeP %s %s.zip %s", password, outputName, payloadName)
		}
	} else if method == "tar" {
		finalOutputName = fmt.Sprintf("%s.tar", outputName)
		command = fmt.Sprintf("tar -cvf %s %s", finalOutputName, payloadName)
	} else if method == "tar.gz" {
		finalOutputName = fmt.Sprintf("%s.tar.gz", outputName)
		command = fmt.Sprintf("tar -czvf %s %s", finalOutputName, payloadName)
	} else if method == "tar.bz2" {
		finalOutputName = fmt.Sprintf("%s.tar.bz2", outputName)
		command = fmt.Sprintf("tar -cjvf %s %s", finalOutputName, payloadName)
	} else if method == "tar.xz" {
		finalOutputName = fmt.Sprintf("%s.tar.xz", outputName)
		command = fmt.Sprintf("tar -cJvf %s %s", finalOutputName, payloadName)
	}

	// Setup a temporary directory to build the payload in
	agent_build_path, err := os.MkdirTemp("", "compressor-wrapper")
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = fmt.Sprintf("Error creating a temp directory: %v", err)
		customLogger.Printf("Error creating a temp directory: %v\n", err)
		return payloadBuildResponse
	}
	defer os.RemoveAll(agent_build_path)

	// Write the payload to disk
	err = os.WriteFile(filepath.Join(agent_build_path, payloadName), *payloadBuildMsg.WrappedPayload, os.FileMode(0644))
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = fmt.Sprintf("Error writing %s to disk: %v", payloadName, err)
		customLogger.Printf("Error writing %s to disk: %v\n", payloadName, err)
		return payloadBuildResponse
	}

	// Run the command to compress the payload
	_, err = runCommand(command, agent_build_path)
	if err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildStdErr = fmt.Sprintf("Error running command (%s): %v", command, err)
		customLogger.Printf("Error running command (%s): %v\n", command, err)
		return payloadBuildResponse
	}

	// Read the final output to provide back to the user
	finalFileLocation := fmt.Sprintf("%s/%s", agent_build_path, finalOutputName)
	if fileBytes, err := os.ReadFile(finalFileLocation); err != nil {
		payloadBuildResponse.Success = false
		payloadBuildResponse.BuildMessage = "Failed to find final file"
	} else {
		payloadBuildResponse.Payload = &fileBytes
		payloadBuildResponse.Success = true
		payloadBuildResponse.BuildMessage = "Successfully compressed the payload!"
	}

	customLogger.Println("Successfully completed build process")

	// Return the payload or archived payload
	return payloadBuildResponse
}

type CommandOutput struct {
	Stdout string
	Stderr string
}

// Run OS command
func runCommand(command string, directory string) (CommandOutput, error) {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(command)
	cmd.Dir = directory
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := CommandOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	return output, err
}

func main() {
	agentstructs.AllPayloadData.Get("compressor").AddPayloadDefinition(compressorDefinition)
	agentstructs.AllPayloadData.Get("compressor").AddBuildFunction(build)
	agentstructs.AllPayloadData.Get("compressor").AddIcon("compressor.svg")

	MythicContainer.StartAndRunForever([]MythicContainer.MythicServices{
		MythicContainer.MythicServicePayload,
	})
}
