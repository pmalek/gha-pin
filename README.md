# gha-pin

Pin Github Actions versions with commit SHA

```diff
@@ -13,15 +13,15 @@ jobs:
   unit-tests:
     runs-on: ubuntu-latest
     steps:
-    - uses: actions/checkout@v4.2.2
+    - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
     - run: make test.unit
 
   integration-tests:
     runs-on: ubuntu-latest
     steps:
-    - uses: actions/checkout@v4.2.2
+    - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
     - run: make test.unit
     - name: setup golang
-      uses: actions/setup-go@v5.2.0
+      uses: actions/setup-go@3041bf56c941b39c61721a86cd11f3bb1338122a # v5.2.0
     - name: run integration tests
       run: make test.integration
```

## How to install

```
go install github.com/pmalek/gha-pin/cmd/gha-pin@latest
```

## How to use

```
gha-pin <path-to-workflow-file> ...
```

Optionally provide `GITHUB_TOKEN` environment variable to increase rate limit for Github API.
