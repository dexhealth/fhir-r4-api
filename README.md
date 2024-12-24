# FHIR R4 HTTP REST API

An HTTP REST API server implementing the HL7 FHIR R4 specification.

## How to use

```shell
go run ./cmd
```

## Import Private Repo

When the dependent package (dexhealth/fhir) changes, we will need update the `go.mod` to reflect the new version.

### fatal: could not read Username

Problem

```shell
go mod tidy

go: downloading github.com/dexhealth/fhir v0.0.0-20241221200152-a2fccf9a969b
go: github.com/dexhealth/fhir@v0.0.0-20241221200152-a2fccf9a969b: verifying module: github.com/dexhealth/fhir@v0.0.0-20241221200152-a2fccf9a969b: reading https://sum.golang.org/lookup/github.com/dexhealth/fhir@v0.0.0-20241221200152-a2fccf9a969b: 404 Not Found
        server response:
        not found: github.com/dexhealth/fhir@v0.0.0-20241221200152-a2fccf9a969b: invalid version: git ls-remote -q origin in /tmp/gopath/pkg/mod/cache/vcs/a990e72f6f7aa215cb56194706241f7af858a545e1a26ede3fd2a8b76e0671e6: exit status 128:
                fatal: could not read Username for 'https://github.com': terminal prompts disabled
        Confirm the import path was entered correctly.
        If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
```

Solution

```sh
git config --global url.ssh://git@github.com/.insteadOf https://github.com/
GOPRIVATE=github.com/dexhealth/fhir
go get github.com/dexhealth/fhir
```
