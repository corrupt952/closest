# closest

[![Go Report Card](https://goreportcard.com/badge/github.com/corrupt952/closest)](https://goreportcard.com/report/github.com/corrupt952/closest)
[![Test](https://github.com/corrupt952/closest/actions/workflows/test.yml/badge.svg)](https://github.com/corrupt952/closest/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A lightweight command-line tool that searches the current directory or parent directories for specific files and returns the closest path.

## Why closest?

Many tools only look for configuration files in the current directory or in the home directory, but not in parent directories. This makes it difficult to use these tools in monorepos or nested project structures.

`closest` solves this problem by finding the nearest matching file in the directory hierarchy, making it easy to:

- Use tools with configuration files in parent directories
- Find configuration files in monorepo structures
- Locate files without knowing their exact location
- Troubleshoot configuration inheritance

## Installation

### Download binary

Download the latest binary from [GitHub Releases](https://github.com/corrupt952/closest/releases).

```sh
# Example for Linux (amd64)
curl -L https://github.com/corrupt952/closest/releases/latest/download/closest_linux_amd64.tar.gz | tar xz
sudo mv closest /usr/local/bin/
```

### Install via aqua

If you use [aqua](https://github.com/aquaproj/aqua), you can install `closest` with:

```sh
aqua g -i corrupt952/closest
```

### Build from source

```sh
git clone https://github.com/corrupt952/closest.git
cd closest
go build
```

## Usage

```sh
Usage: closest [options] [pattern]
Options:
  -a    Search all files[default: false]
  -r    Use regex pattern for matching[default: false]
  -v    Show version
```

### Exit Codes

The tool uses the following exit codes:

- `0`: Success - Files were found and output
- `1`: Error - An error occurred (file not found, invalid regex, permission denied, etc.)

### Basic Usage

Find the closest file matching a specific name:

```sh
closest .tflint.hcl
# Output: /path/to/closest/.tflint.hcl
```

Find all matching files from current directory to root:

```sh
closest -a .envrc
# Output: 
# /current/path/.envrc
# /current/.envrc
# /home/user/.envrc
```

Find files using regex patterns:

```sh
closest -r ".*\.ya?ml$"
# Output: /path/to/closest/config.yaml
```

## Examples

### Example 1: Using with tflint

`tflint` only references `.tflint.hcl` in the current or home directory. With `closest`, you can use project-specific settings from parent directories in a monorepo.

Directory structure:
```
/
└── home
    └── app
        └── terraform
            ├── .tflint.hcl
            └── example-service
                ├── production
                └── staging # <- current directory
```

Run tflint with the closest configuration:

```sh
tflint --config $(closest .tflint.hcl)
```

### Example 2: Troubleshooting direnv configuration

When using `direnv`, you might want to find all `.envrc` files that affect the current directory. The `-a` option helps with troubleshooting by showing all relevant files.

Directory structure:
```
/
└── home
    └── app
        ├── .envrc
        └── terraform
            ├── .envrc
            └── example-service
                ├── .envrc
                ├── production  # <- current directory
                |   └── .envrc
                └── staging
```

Find all `.envrc` files from current directory to root:

```sh
closest -a .envrc
```

Output:
```sh
/home/app/terraform/example-service/production/.envrc
/home/app/terraform/example-service/.envrc
/home/app/terraform/.envrc
/home/app/.envrc
```

> **Note:** Command options must come before the filename. For example, `closest .envrc -a` doesn't work.

## Error Handling

`closest` provides clear error messages for common issues:

- **File not found**: When no matching files are found in the directory hierarchy
  ```
  Error: file not found: config.json
  ```

- **Invalid regex pattern**: When the provided regex pattern is invalid
  ```
  Error: invalid regex pattern: error parsing regexp: missing closing ]: `[abc`
  ```

- **Permission denied**: When the tool cannot access a directory due to permission issues
  ```
  Error: failed to read directory /path/to/restricted: permission denied
  ```

- **Missing pattern argument**: When no search pattern is provided
  ```
  Error: error parsing flags: missing pattern argument
  ```

### Example 3: Finding configuration files with regex

Sometimes you need to find configuration files that might have different extensions. The `-r` option enables regex pattern matching for flexible searches.

Find the closest YAML configuration file:

```sh
closest -r ".*\.ya?ml$"
```

Find all YAML files in the directory hierarchy:

```sh
closest -a -r ".*\.ya?ml$"
```

Output:
```sh
/home/app/terraform/example-service/production/config.yaml
/home/app/terraform/example-service/terraform.yaml
/home/app/terraform/main.yml
/home/app/config.yml
```

## Troubleshooting

### Common Issues

1. **No files found**
   - Check if the file exists in any parent directory
   - Verify the spelling and case of the filename (filesystems may be case-sensitive)
   - If using regex, ensure the pattern is correct

2. **Permission errors**
   - Ensure you have read permissions for all directories in the path
   - Try running with elevated privileges if necessary

3. **Regex not matching**
   - Test your regex pattern with a regex testing tool
   - Remember to escape special characters
   - For complex patterns, start simple and build up

### Debugging Tips

- Use the `-a` flag to see all matching files, which can help identify if files exist but aren't where expected
- For regex issues, try simplifying your pattern first, then make it more specific

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

For detailed information about the development setup, workflow, and release process, please see the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
