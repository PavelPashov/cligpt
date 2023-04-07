# CLIGPT

CLIGPT is a command-line interface (CLI) application for interacting with OpenAI's ChatGPT. The application allows users to start a chat with the model, change the model configuration, prompt the model with a single prompt, and update the token.

## Demo

### Chat

![](https://github.com/EiTamOnya/cligpt/examples/chat.gif)

### List previous chat sessions

![](https://github.com/EiTamOnya/cligpt/examples/list.gif)

## Building the App

To install CLIGPT, follow these steps:

1. Clone the repository: `git clone https://github.com/EiTamOnya/cligpt.git`
2. Navigate to the project directory: `cd cligpt`
3. Run the build script: `./build cligpt`

The build script supports different architectures and platforms. You can modify the `platforms` array in the `build.sh` file to include the architectures you want to build for. For example, the default `platforms` array includes the following architectures:

```
platforms=("windows/amd64" "linux/amd64" "darwin/amd64")
```

This will build the application for Windows, Linux, and macOS on 64-bit architectures.

## Installation

### Linux

To install CLIGPT on Linux, build or download `cligpt-linux-amd64` from [here](https://github.com/EiTamOnya/cligpt/releases/latest).

Then, run the following command:

```
sudo cp <file> /usr/bin/cligpt
```

Replace `<file>` with the name of the downloaded file.

### macOS

To install CLIGPT on macOS, build or download `cligpt-darwin-amd64` from [here](https://github.com/EiTamOnya/cligpt/releases/latest).

Then, run the following command:

```
sudo cp <file> /usr/local/bin/cligpt
```

Replace `<file>` with the name of the downloaded file.

## Configuration

Before using CLIGPT, you will need to create an OpenAI API key. To do this, follow these steps:

1. Sign up for an OpenAI API key [here](https://beta.openai.com/signup/).
2. Create a new API key (you might need to add a payment method if you're not eligible for a free trial).
3. Run `cligpt init` and copy the key once prompted.

## Available Commands

To use CLIGPT, run the following commands:

- `cligpt chat`: Start a chat with the model.
- `cligpt init`: Initiate the setup for CLIGPT.
- `cligpt model`: Change the model configuration.
- `cligpt prompt`: Prompt the model with a single prompt.
- `cligpt token`: Update the token.
- `cligpt persona`: Select a personality for the model.
- `cligpt maxt`: Set the number of max tokens.
- `cligpt temp`: Set the temperature.

## Contributing

If you would like to contribute to CLIGPT, please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b my-feature-branch`
3. Make your changes and commit them: `git commit -m "Add my feature"`
4. Push your changes to your fork: `git push origin my-feature-branch`
5. Create a pull request.
