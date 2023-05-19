# closest

The command that searches the current directory or parent directories for a specific file and returns the closest path

## Usage

```sh
Usage: closest [options] [pattern]
Options:
  -a    Search all files[default: false]
```

To find a closest file the current directory, run the following:

```sh
closest .tflint.hcl
```

To find all files from the current directory to root directory, run the following:

```sh
closest -a .envrc
```

### Example 1: Find a .tflint.hcl file and run tflint


`tflint` only references `.tflint.hcl` in the current or home directory.
This makes it easy to read per-project settings in the repository root or in the terraform directory in monorepo.

The directory structure is as follows, where `staging` is the current directory.

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

To run `tflint` in combination with `closest`, run the following:

```sh
tflint --config $(closest .tflint.hcl)
```

### Example 2: Find all .envrc files up to the root directory

Sometimes when using `direnv`, you want to find where `.envrc` is defined from the root directory to the current directory.
In that case, you can use the `-a` option to display all `.envrc` files up to the root directory, which is useful for troubleshooting.

The directory structure is as follows, where `production` is the current directory.

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

To find all `.envrc` from `production` to the root directory, run the following:

```sh
closest -a .envrc
```

Please take care that **the filename must be prefix with `-a`**.
For example, `closest .envrc -a` doesn't work.

The output:

```sh
/home/app/terraform/example-service/production/.envrc
/home/app/terraform/example-service/.envrc
/home/app/terraform/.envrc
/home/app/.envrc
```
