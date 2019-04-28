# README Documentation

# Table of Contents

- [Downloading, Installing, and Running Our Product](#downloading-installing-and-running-our-product)
  - [Downloading](#downloading)
    - [From GitHub](#from-github)
    - [From `master.zip`](#from-masterzip)
    - [The ISO Image of the Books](#the-iso-image-of-the-books)
  - [Installing](#installing)
    - [Natural Language Toolkit (NLTK)](#natural-language-toolkit-nltk)
    - [PySpark Dependencies](#pyspark-dependencies)
    - [Golang](#golang)
    - [Gorilla WebSocket](#gorilla-websocket)
    - [Our Project](#our-project)
    - [The ISO Image of the Books](#the-iso-image-of-the-books-1)
    - [The HAR Archive from the Books](#the-har-archive-from-the-books)
  - [Running](#running)
    - [Data Cleaning](#data-cleaning)
    - [Main Program](#main-program)
    - [Streaming Version of Main Program](#streaming-version-of-main-program)
    - [Web Interface](#web-interface)
- [File Listings](#file-listings)

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

Please type the following command to install the library NLTK.

```sh
sudo pip install -U nltk
```

Punkt Sentence Tokenizer divides a text into a list of sentences by using an unsupervised algorithm to build a model for abbreviation words, collocations, and words that start sentences. 

To download "punkt" tokenizer, run the following command in your python shell.
```sh
import nltk
nltk.download('punkt')
```

### PySpark Dependencies

To run the project with the required dependency, you need to create a file named `requirements.txt` in the same directory of the python program.
Add below package into `requirements.txt` and save.
```sh
nltk==3.4.1
pyspark==2.4.0
```
Run below command in your console. Use apt-get to install zip if you don't have it installed.
```sh
pip install -t dependencies -r requirements.txt
cd dependencies
zip -r ../dependencies.zip 
```

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

### The HAR Archive from the Books

Please run the data cleaning step once (see below) before running this step.

To create the Hadoop Archive Format for grouping the small ebooks files, the following command should be used:

```sh
$ hadoop archive -archiveName booksarchive.har -p /books /har
```
where the source ebooks files are placed under `hdfs://books/` 
and the archive created will be stored as `hdfs://har/booksarchive.har`



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
spark-submit --master yarn m4.py <mode> <books_path> <mystery_text> \
    <min_partitions> <preserves_partitioning> <num_partitions>
```

#### Parameters

| Parameter | Description |
| --- | --- |
| `mode` | `0`: Find the result without reading the cache. Do not save the result to the cache.<br>`1`: Find the result without reading the cache. Save the result to the cache.<br>`2`: Do not find the result. Load the result from the cache. |
| `books_path` | The path to the books. You may fill in `har:///har/booksarchive.har` for the HAR version of the book archive, or `hdfs://gpu1:8020/books` for the normal HDFS version of the books. |
| `mystery_text` | The mystery text (initialism) to be searched from the books. |
| `min_partitions` | The `MinPartitions` parameter when calling `sc.wholeTextFiles()`. This has significant effect on the number of partitions created from the source text files. |
| `preserves_partitioning` | The `preservesPartitioning` parameter when calling `books.flatMap()`. This has just little effect on the efficiency, perhaps because Spark already preserves partitioning in our case, without explicitly specified. |
| `num_partitions` | The `numPartitions` parameter when calling `filtered_list_of_book_name_and_sentence.groupByKey()`. If this values is `0`, then the default is used (as if `None` is passed to the function). |


### Streaming Version of Main Program

To run the streaming version of main program from the command line, please change the current directory to `app`, then execute `spark-submit`.

```sh
cd $GOPATH/src/github.com/comp7305-group-1/group-project/app/
spark-submit --master yarn stream2.py 
```

| Parameter | Description |
| --- | --- |
| `books_path` | The path to the books. Please create this directory and leave it empty when starting the program. After starting, you may copy files to this directory. The program will monitor this directory and run the search on-the-fly. |
| `mystery_text` | The mystery text (initialism) to be searched from the books. |
| `batch_duration` | The interval (in seconds) between successive scans. For example, if this value is `1`, then the program checks whether a new file is copied to `books_path` every second. A value of `8` is recommended. |
| `preserves_partitioning` | The `preservesPartitioning` parameter when calling `books.flatMap()`. This has just little effect on the efficiency, perhaps because Spark already preserves partitioning in our case, without explicitly specified. |

### Web Interface

To run the streaming version of main program from the command line, please change the current directory to `web`, then execute `go`.

```sh
go run main.go
```

To use the web interface, please visit port `12345` of the server via `http`.

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





# File Listings
