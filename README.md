# Go stress test CLI

## Running

You can run the latest version of this application using `Docker`:
```sh
docker run eliasfeijo/go-stress-test --url https://google.com --requests 10 --concurrency 3
```

You can also install the application with `Go`:
```sh
go install github.com/eliasfeijo/go-stress-test@latest
# Run the application
go-stress-test --url https://google.com --requests 10 --concurrency 3
```

If you don't have `Go` installed in your system, you can download the application binary from the [latest Release](https://github.com/eliasfeijo/go-stress-test/releases/latest) for your specific OS and architecture, and then extract and execute the binary:
```sh
# Run the application (change the "~/Downloads/" path to where you extracted the binary)
~/Downloads/go-stress-test --url https://google.com --requests 10 --concurrency 3
```
