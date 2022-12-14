# closest

The command that searches the current directory or parent directories for a specific file and returns the closest path

## Usage


## Example

The directory structure is as follows, where `staging` is the current directory.

```
/
└── home
    └── john
        └── terraform
            ├── .tflint.hcl
            └── example-service
                ├── production
                └── staging # <- current directory
```

If you want to search for a .tflint.hcl, you can run `closest .tflint.hcl`, which will return `/home/john/terraform/.tflint.hcl`.

```sh
tflint --config $(closest .tflint.hcl)
```
