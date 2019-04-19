import os
import zipfile
import sys
import nltk
from pyspark import SparkContext, SparkConf
from pyspark.sql import SparkSession
from pyspark.streaming import StreamingContext


def remove_existing_file(path):
    URI = sc._gateway.jvm.java.net.URI
    Path = sc._gateway.jvm.org.apache.hadoop.fs.Path
    FileSystem = sc._gateway.jvm.org.apache.hadoop.fs.FileSystem
    fs = FileSystem.get(URI("hdfs://gpu1:8020"), sc._jsc.hadoopConfiguration())
    fs.delete(Path(path))


def initialism(sentence):
    words = sentence.split()
    initialism = []
    for word in words:
        initialism.append(word[0])
    return ''.join(initialism).lower()


def check(sentence):
    print("check\n")
    if user_input.lower() in initialism(sentence):
        print("Possible mystery found: \n%s" % (sentence))
        return (sentence)
    else:
        return ("")


def get_cleaned_sentence_from_zip_txt(book):
    sentences = nltk.tokenize.sent_tokenize(book[1].replace('\r\n', " ").replace('\n', " "))
    return {"book_name": book[0], "sentences": sentences}


def check_each_book(processed_book):
    list_of_match = []
    book_name = processed_book['book_name']
    sentences = processed_book['sentences']
    for sentence in sentences:  # sentences = sentences list in a book
        if user_input.lower() in initialism(sentence):
            list_of_match.append(sentence.encode('ascii', 'ignore'))
    if len(list_of_match) > 0:
        return {"book_name": book_name, "sentences": list_of_match}


# Read here
def is_exist():
    try:
        print("*************************\n")
        print(sc.textFile("hdfs://gpu1:8020/results/%s" % user_input).collect())
        print("\nResult exist in database\n")
        print("*************************\n")
        return True
    except:
        # get error when result no exist
        return False
    return (is_exist)


def start():
    print("start\n")
    print("Start Mystery===============================================\n")
    # 1. Get sentences from each book
    books = sc.wholeTextFiles('hdfs://gpu1:8020/books')
    processed_books = books.map(get_cleaned_sentence_from_zip_txt)  # [ [s1, s2,s3] [] [] ]
    # 2. Get match result from each book
    result = processed_books.map(check_each_book).filter(lambda x: x is not None)
    print("*************************\n")
    print("Result\n")
    print(result.collect())
    print("*************************\n")
    # 3. Save the result to hdfs for later usage
    result_path = "hdfs://gpu1:8020/results/%s" % user_input
    try:
        # Save the result
        result.saveAsTextFile("hdfs://gpu1:8020/results/%s" % user_input)
    except:
        # Remove the old result file
        remove_existing_file(result_path)
        result.saveAsTextFile("hdfs://gpu1:8020/results/%s" % user_input)


if __name__ == '__main__':
    conf = SparkConf().setAppName("app")
    sc = SparkContext(conf=conf)
    sc.setLogLevel("ERROR")
    user_input = str(sys.argv[1]).lower()
    mode = int(sys.argv[2])
    # default mode 1 = generate the result again, mode 2 = check exist
    if mode == 1 or not is_exist():
        start()
    else:
        pass
    print("End Mystery===============================================\n")