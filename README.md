# cligpt

Cligpt is a command-line interface (CLI) application for interacting with OpenAI's ChatGPT. The application allows users to start a chat with the model, change the model configuration, save conversations and continue them.

## Demo

### Chat

![](https://github.com/EiTamOnya/cligpt/blob/master/examples/chat.gif)

### List previous chat sessions

![](https://github.com/EiTamOnya/cligpt/blob/master/examples/list.gif)

## Building the App

To install cligpt, follow these steps:

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

To install cligpt on Linux, build or download `<file>` from [here](https://github.com/EiTamOnya/cligpt/releases/latest).

Then, run the following command:

```
sudo cp <file> /usr/bin/cligpt
```

Replace `<file>` with the name of the downloaded file.

### macOS

To install cligpt on macOS, build or download `<file>` from [here](https://github.com/EiTamOnya/cligpt/releases/latest).

Then, run the following command:

```
sudo cp <file> /usr/local/bin/cligpt
```

Replace `<file>` with the name of the downloaded file.

## Configuration

Before using cligpt, you will need to create an OpenAI API key. To do this, follow these steps:

1. Sign up for an OpenAI API key [here](https://beta.openai.com/signup/).
2. Create a new API key (you might need to add a payment method if you're not eligible for a free trial).
3. Run `cligpt init` and add the key once prompted.

## Available Commands

These are the available commands for cligpt:

- `cligpt chat`: Start a chat with the model.
- `cligpt init`: Initiate the setup for cligpt.
- `cligpt model`: Select a model which will be saved to your config.
- `cligpt prompt`: Prompt the model with a single prompt.
- `cligpt token`: Update the your OpenAI API key.
- `cligpt persona`: Select a personality for the model. This is used in the first system message if provided.
- `cligpt maxt`: Set the number of max tokens to generate in the chat completion.
- `cligpt temp`: Set the sampling temperature.

Use `--help` or `-h` after any command to see the available subcommands and prompts.

## Contributing

If you would like to contribute to cligpt, please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b my-feature-branch`
3. Make your changes and commit them: `git commit -m "Add my feature"`
4. Push your changes to your fork: `git push origin my-feature-branch`
5. Create a pull request.
