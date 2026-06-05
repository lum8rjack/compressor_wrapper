# Compressor Wrapper

A wrapper payload for [Mythic](https://github.com/its-a-feature/Mythic) that wraps any binary (any file type or OS) and packages it into the specified archive format. This makes it easier to package the generated payload via Mythic without needing to download the agent and do it manually. An examples is adding the payload to a password protected zip file.

## Supported Methods

Examples commands that are ran for each of the methods:

```bash
# Normal zip
zip -j package.zip agent.exe

# Zip with a password
zip -jeP password package.zip agent.exe

# tar
tar -cvf package.tar agent.exe

# tar (Gzip)
tar -czvf package.tar.gz agent.exe

# tar (Bzip2)
tar -cjvf package.tar.bz2 agent.exe

# tar (XZ)
tar -cJvf package.tar.xz agent.exe
```

## Agent Setup

You will need to edit the payload code to allow it to be wrapperd by compressor. Below is the [Poseidon](https://github.com/MythicAgents/poseidon/blob/master/Payload_Type/poseidon/poseidon/agentfunctions/builder.go) code updated. You need to make sure the "CanBeWrappedByTheFollowingPayloadTypes" includes "compressor".

```go
// Payload_Type/poseidon/poseidon/agentfunctions/builder.go

var payloadDefinition = agentstructs.PayloadType{
	Name:                                   "poseidon",
	SemVer:                                 version,
	FileExtension:                          "bin",
	Author:                                 "@xorrior, @djhohnstein, @Ne0nd0g, @its_a_feature_",
	SupportedOS:                            []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_MACOS},
	Wrapper:                                false,
	CanBeWrappedByTheFollowingPayloadTypes: []string{"compressor"},
```

## Usage

You can install the wrapper the same way you install any other agent by using the `mythic-cli` command.

```bash
./mythic-cli install github https://github.com/lum8rjack/compressor_wrapper
```

Once installed, you can log into Mythic and wrap a payload using the following steps:

- Click **Create Wrapper** on the left sidebar
    - Select the operating system
    - Select **compressor** for the payload type
    - Select **Start Fresh**
- On the next screen complete each field
    - Name: Specify the file name of the payload inside the compressed archive
    - Method: Specify the method to use to compress the payload
    - If you selected "Zip", there will be an option to add a password
- On the last screen select the payload you aleady build and want to compress
- Then click **Create Payload**

Once the generation is complete, go to **Payloads** to download the output file.
