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


def initialism(sentence): # sentence: str
    '''
    Extracts the first letter of each of the words in the sentence,
    concatenates them into a string, and turns it into lowercase.
    For instance, if sentence is 'I go to school by bus', then 'igtsbb' is
    returned. If sentence is 'Cloud Computing is so fun!', then 'ccisf' is
    returned.
    '''
    return ''.join([ word[0] for word in sentence.split() ]).lower()


def split_content_into_sentences(book): # book: tuple(unicode, unicode) = (book_name, book_content)
    '''
    Splits the content of the book into sentences. The content is converted
    into ASCII and unconvertible characters are ignored. The parameter is a
    2-element tuple. The first element is the unique name of the book (in
    Unicode). The second element is the content of the book (in Unicode). The
    returned value is a 2-element tuple. The first element is the unique name
    of the book (in ASCII). The second element is a list of sentences from the
    content of the book (in ASCII).
    '''
    book_name = book[0].encode('ascii', 'ignore') # book_name: str
    book_content = book[1].encode('ascii', 'ignore') # book_content: str
    book_content_in_one_line = book_content.replace('\r\n', ' ').replace('\n', ' ') # book_content_in_one_line: str
    list_of_sentences = nltk.tokenize.sent_tokenize(book_content_in_one_line) # list_of_sentences: list = [ sentence1, sentence2, ... ]
    return [ (book_name, sentence) for sentence in list_of_sentences ]


def get_result_path(mystery_text):
    return 'hdfs://gpu1:8020/results/%s' % mystery_text


def find_result(sc, mystery_text, books, min_partitions): # By default, books = 'hdfs://gpu1:8020/books'
    books = sc.wholeTextFiles(books, minPartitions=min_partitions, use_unicode=True) # books: RDD; Error after setting use_unicode=False
    list_of_book_name_and_sentence = books.flatMap(split_content_into_sentences) # list_of_book_name_and_sentence: PipelinedRDD = [ (book_name_1, sentence1), (book_name_1, sentence2), (...), (book_name_2, ...), ...]
    result = list_of_book_name_and_sentence.filter(lambda x: mystery_text in initialism(x[1]))
    collected = result.groupByKey().mapValues(list).collect()
    print('=== Begin Result ===============================================================')
    print(collected)
    print('=== End Result =================================================================\n')
    return result


def save_result(result, mystery_text):
    path = get_result_path(mystery_text)
    try:
        # Save the result
        result.saveAsTextFile(path)
    except:
        # Remove the old result file
        remove_existing_file(path)
        result.saveAsTextFile(path)


def load_result(sc, mystery_text):
    try:
        path = get_result_path(mystery_text)
        collected = sc.textFile(path).collect()
        print('=== Begin Result ===============================================================')
        print(collected)
        print('=== End Result =================================================================\n')
    except:
        print('*** ERROR: Result does not exist in database ***')


if __name__ == '__main__':
    print('=== Begin Mystery ==============================================================\n')
    conf = SparkConf().setAppName('UnravelMysteries').setMaster('yarn')
    sc = SparkContext(conf=conf)
    sc.setLogLevel('INFO')
    # mode 0 = generate the result from scratch
    # mode 1 = generate the result from scratch + save result
    # mode 2 = load result from saved result
    mode = int(sys.argv[1])
    #
    if mode == 0 or mode == 1:
        mystery_text = str(sys.argv[2]).lower()
        books = sys.argv[3]
        min_partitions = int(sys.argv[4])
        result = find_result(sc, mystery_text, books, min_partitions)
        if mode == 1:
            save_result(result, mystery_text)
    elif mode == 2:
        mystery_text = str(sys.argv[2]).lower()
        load_result(sc, mystery_text)
    else:
        print('*** ERROR: Invalid mode ***')
    print('=== End Mystery ================================================================\n')
