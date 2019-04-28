# Downloading, Installing, and Running Our Product


## Downloading

Please download the program from GitHub, or find it in the file `master.zip` we submitted.

### From GitHub

```sh
git clone https://github.com/comp7305-group-1/group-project.git
```

### From `mystery.zip`

```sh
unzip master.zip
```


## Installing

The main program was written in Python, and it depends on the PySpark component and the Natural Language Toolkit (NLTK) library.

The web interface was written in Golang, and it depends on Gorilla WebSocket library.

### Natural Language Toolkit (NLTK)

***TODO: Nelson, please add this!***

### Golang

Please visit https://golang.org/doc/install and follow the instructions to install the Golang distribution. This typically involves just a few steps:

- Downloading the distribution from the URL above
- Decompressing the downloaded archive and placing it at the desired location
- Configuring the environment variables so that the `go` executable is in the `$PATH`
- Creating a directory for Golang projects and setting the environment variable `$GOPATH` to point to it

### Gorilla WebSocket

Please type the following command to install the library under `$GOPATH`.

```sh
go get -u github.com/gorilla/websocket
```

### Our Project

It would be easier to run our project under `$GOPATH`. In the following instructions, let `group-project` be the directory downloaded from GitHub (or unzipped from the archive).

```sh
mkdir -p $GOPATH/src/github.com/comp7305-group-1
mv group-project $GOPATH/src/github.com/comp7305-group-1/
```


## Running

### Data Cleaning

### Main Program

### Streaming Version of Main Program
