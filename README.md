![License](https://img.shields.io/badge/license-MIT-blue)
[![testing](https://github.com/ondrovic/bambulab-authenticator/actions/workflows/testing.yml/badge.svg)](https://github.com/ondrovic/bambulab-authenticator/actions/workflows/testing.yml)
[![releaser](https://github.com/ondrovic/bambulab-authenticator/actions/workflows/releaser.yml/badge.svg)](https://github.com/ondrovic/bambulab-authenticator/actions/workflows/releaser.yml)

# Bambulab Authenticator CLI

The Bambulab Authenticator CLI is a command-line tool that allows you to authenticate with your Bambulab credentials and save the authentication information to a file.

## Installation

To install the Bambulab Authenticator CLI, you can use the following command:

go install github.com/ondrovic/bambulab-authenticator/cmd/cli@latest

## Usage

To authenticate with your Bambulab credentials, use the following command:

cli authenticate --user-account <your-account> --user-password <your-password> --user-region <your-region> --output-path <output-path>

Replace the following placeholders:

- `<your-account>`: Your Bambulab user account.
- `<your-password>`: Your Bambulab user password.
- `<your-region>`: Your Bambulab user region.
- `<output-path>`: The path to save the authentication information.

All of the flags are required.

## Development

To build and run the Bambulab Authenticator CLI locally, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/ondrovic/bambulab-authenticator.git
   ```

2. Navigate to the project directory:

   ```
   cd bambulab-authenticator
   ```

3. Build the CLI:

   ```
   go build -o cli cmd/cli/main.go
   ```

4. Run the CLI:

    ```
   ./cli authenticate --user-account <your-account> --user-password <your-password> --user-region <your-region> --output-path <output-path>
   ```

## Contributing

If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.