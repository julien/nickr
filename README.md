# NickR

Requirements
------------

  + Download [Go](http://golang.org/dl)

  + Setup your Go environment

  + Install the Go dependencies

    ```shell
    go get github.com/tools/godep
    go get github.com/melvinmt/firebase
    ```

Development
-----------

  + In a terminal, issue the following command to run the server

    ```shell
    go run main.go
    ```

  + To run tests, do

    ```shell
    go test --coverage ./...
    ```

  + If add new files don't forget to run

    ````shell
    godep save -r ./...
    ```







