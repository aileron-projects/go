# Realease manual

This document shows how to release a new version.

## 1. Checkout release branch

First, clone `main` branch and checkout a new release branch.

Use `release-vX.Y.Z` branch.
The release branch should be deleted after merged into the main branch.

```bash
git checkout -b release-vX.Y.Z
```

Following steps should be done on the release branch.

## 2. Update minimum Go version

Update go minimum version if necessary.

This step is basically required after official Go majour release.

This project supports current latest 3 majour releases.

- `go 1.(N).x`
- `go 1.(N-1).x`
- `go 1.(N-2).x`

To update the minimum Go version, update `go` directive in the `go mod`.

For example, when minimum version is `go1.23.0`, the go directive should be

```text
go 1.23.0
```

## 3. Update dependency

Update dependencies.
Make sure to use `GOTOOLCHAIN` directive with minimum Go version to avoid automatic modification of go.mod file.

When the minimum Go version is go1.23.0, the command will be

```bash
GOTOOLCHAIN=go1.23.0 go get -u ./...
GOTOOLCHAIN=go1.23.0 go mod tidy
```

Note that the command can still raise some errors because of the version incompatibility such as

```bash
$ GOTOOLCHAIN=go1.23.0 go mod tidy
go: github.com/open-policy-agent/opa@v1.2.0 requires go >= 1.23.6 (running go 1.23.0; GOTOOLCHAIN=go1.23.0)
```

In such cases, downgrade the dependencies until all errors gone.

## 4. Push and run tests

After updating dependencies commit and push changes.
The release branch can be merged into the main branch when all workflows have passed.

Fix errors if some errors occurred.

**Optional tasks:**

Update other resources such as workflows when it should be before merging into the main branch.

## 5. Create tag

After passign all tests, merge the release branch into the main branch.
And then, create the release tag `vX.Y.Z` at main branch.
