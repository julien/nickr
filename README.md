# NickR

Requirements
------------

  + Download [Go](http://golang.org/dl)

  + Setup and [test your Go installation](https://golang.org/doc/install#testing)

  + Install the Go dependencies

    ```shell
    go get github.com/tools/godep
    go get github.com/melvinmt/firebase
    ```

Development
-----------

  + In a terminal, issue the following command to run the server

    ```shell
    go build -o main
    ```

    This will compile the program as "main"

    You can then run the program (i.e. launch the server) with:

    ```shell
    ./main
    ```

  + To run tests, do

    ```shell
    go test --coverage ./...
    ```

    If you want to generate test coverage output, use:

    ````shell
    go tests --coverprofile=out.cov ./...
    ```

    You can then view the coverage report in your browser with:

    ```shell
    go tool cover -html=out.cov
    ```

  + If add new files don't forget to run

    ````shell
    godep save -r ./...
    ```


