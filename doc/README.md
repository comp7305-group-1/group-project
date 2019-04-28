# README Documentation

# Table of Contents

- [Downloading, Installing, and Running Our Product](#downloading-installing-and-running-our-product)
  - [Downloading](#downloading)
    - [From GitHub](#from-github)
    - [From `master.zip`](#from-masterzip)
    - [The ISO Image of the Books](#the-iso-image-of-the-books)
  - [Installing](#installing)
    - [Natural Language Toolkit (NLTK)](#natural-language-toolkit-nltk)
    - [Golang](#golang)
    - [Gorilla WebSocket](#gorilla-websocket)
    - [Our Project](#our-project)
    - [The ISO Image of the Books](#the-iso-image-of-the-books-1)
  - [Running](#running)
    - [Data Cleaning](#data-cleaning)
    - [Main Program](#main-program)
    - [Streaming Version of Main Program](#streaming-version-of-main-program)
    - [Web Interface](#web-interface)


# Downloading, Installing, and Running Our Product


## Downloading

Please download the project from GitHub, or find it in the file `master.zip` we submitted.

### From GitHub

```sh
git clone https://github.com/comp7305-group-1/group-project.git
```

### From `master.zip`

```sh
unzip master.zip
```

### The ISO Image of the Books

The ISO image of the books could be downloaded from ftp://mirrors.pglaf.org/mirrors/gutenberg-iso/pgdvd042010.iso.

```sh
wget ftp://mirrors.pglaf.org/mirrors/gutenberg-iso/pgdvd042010.iso
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

### The ISO Image of the Books

To install the ISO image of the books, just copy all the contents inside the ISO image to HDFS.

As `sudo` user (i.e. `student`):

```sh
sudo mkdir mountpoint
sudo mount pgdvd042010.iso mountpoint
```

As Hadoop user (i.e. `hduser`):

```sh
hdfs dfs -put mountpoint /data
```

As `sudo` user (i.e. `student`):

```sh
sudo umount mountpoint
```


## Running

### Data Cleaning

To run the data cleaning, please change the current directory to `cleaning`, then execute `spark-submit`.

```sh
cd $GOPATH/src/github.com/comp7305-group-1/group-project/cleaning/
spark-submit --master yarn clean.py <hadoopMasterIP> <hadoopMasterName> <sparkMasterURL>
```

***TODO: Pauline, please add the descriptions of the parameters!***

### Main Program

To run the main program from the command line, please change the current directory to `app`, then execute `spark-submit`.

```sh
cd $GOPATH/src/github.com/comp7305-group-1/group-project/app/
spark-submit --master yarn m4.py
```

***TODO: Params***

### Streaming Version of Main Program

To run the streaming version of main program from the command line, please change the current directory to `app`, then execute `spark-submit`.

```sh
cd $GOPATH/src/github.com/comp7305-group-1/group-project/app/
spark-submit --master yarn stream2.py 
```

***TODO: Params***

### Web Interface

To run the streaming version of main program from the command line, please change the current directory to `web`, then execute `go`.

```sh
go run main.go
```

To use the web interface, please visit port `12345` of the server.

To stop the web server, please kill the process by pressing `Ctrl-C`.

You may also want to run the web interface in background:

```sh
go run main.go &
disown %1
```

To stop the web server, please kill the process by using the `kill` command.

The process ID of the process could be found by the `ps` command:

```sh
ps -A -opid,ppid,user,cmd= | grep go
```

A typical output would be something similar to this:

```
 7586  4092 hduser   grep --color=auto go
 7729     1 hduser   go run main.go
 7765  7729 hduser   /tmp/go-build563353986/b001/exe/main
```

Please kill both `7729` and `7765` in this case:

```sh
kill 7729 7765
```
