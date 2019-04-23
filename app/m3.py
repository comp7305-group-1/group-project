import nltk
import os
import sys
import zipfile
from pyspark import SparkContext, SparkConf
from pyspark.sql import SparkSession
from pyspark.streaming import StreamingContext


def remove_existing_file(path):
    URI = sc._gateway.jvm.java.net.URI
    Path = sc._gateway.jvm.org.apache.hadoop.fs.Path
    FileSystem = sc._gateway.jvm.org.apache.hadoop.fs.FileSystem
    fs = FileSystem.get(URI('hdfs://gpu1:8020'), sc._jsc.hadoopConfiguration())
    fs.delete(Path(path))


def initialism(sentence):
    return ''.join([ word[0] for word in sentence.split() ]).lower()


def get_cleaned_sentence_from_zip_txt(book): # book: tuple(unicode, unicode), book[0]: book_name, book[1]: book_content
    sentences = nltk.tokenize.sent_tokenize(book[1].replace('\r\n', ' ').replace('\n', ' '))
    return {'book_name': book[0], 'sentences': sentences}


def check_each_book(processed_book):
    list_of_match = []
    book_name = processed_book['book_name']
    sentences = processed_book['sentences']
    for sentence in sentences:  # sentences = sentences list in a book
        if user_input.lower() in initialism(sentence):
            list_of_match.append(sentence.encode('ascii', 'ignore'))
    if len(list_of_match) > 0:
        return {'book_name': book_name, 'sentences': list_of_match}


def get_result_path(user_input):
    return 'hdfs://gpu1:8020/results/%s' % user_input


def find_result(sc, user_input, min_partitions):
    # 1. Get sentences from each book
    books = sc.wholeTextFiles('hdfs://gpu1:8020/books', minPartitions=min_partitions, use_unicode=True) # books: RDD; Error after setting use_unicode=False
    processed_books = books.map(get_cleaned_sentence_from_zip_txt) # processed_books: PipelinedRDD; [ (b1, [s1, s2, s3]), (), () ]
    # 2. Get match result from each book
    result = processed_books.map(check_each_book).filter(lambda x: x is not None)
    collected = result.collect()
    print('=== Start Result ===============================================================')
    print(collected)
    print('=== End Result =================================================================\n')
    # 3. Save the result to hdfs for later usage
    path = get_result_path(user_input)
    #try:
    #    # Save the result
    #    result.saveAsTextFile(path)
    #except:
    #    # Remove the old result file
    #    remove_existing_file(path)
    #    result.saveAsTextFile(path)


def load_result(sc, user_input):
    try:
        path = get_result_path(user_input)
        collected = sc.textFile(path).collect()
        print('=== Start Result ===============================================================')
        print(collected)
        print('=== End Result =================================================================\n')
    except:
        print('*** ERROR: Result does not exist in database ***')


if __name__ == '__main__':
    print('=== Start Mystery ==============================================================\n')
    conf = SparkConf().setAppName('UnravelMysteries').setMaster('yarn')
    sc = SparkContext(conf=conf)
    sc.setLogLevel('INFO')
    user_input = str(sys.argv[1]).lower()
    mode = int(sys.argv[2])
    # default mode 1 = generate the result again, mode 2 = check exist
    if mode == 1:
        min_partitions = int(sys.argv[3])
        find_result(sc, user_input, min_partitions)
    elif mode == 2:
        load_result(sc, user_input)
    else:
        print('*** ERROR: Invalid mode ***')
    print('=== End Mystery ================================================================\n')
