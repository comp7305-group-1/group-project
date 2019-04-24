import nltk
import os
import sys
import zipfile
from pyspark import SparkContext, SparkConf


def remove_existing_file(path):
    '''
    Deletes an existing file from HDFS.
    '''
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
    list_of_sentences = nltk.tokenize.sent_tokenize(book_content_in_one_line) # list_of_sentences: list = [ sentence_1, sentence_2, ... ]
    return [ (book_name, sentence) for sentence in list_of_sentences ]


def get_result_path(mystery_text):
    '''
    Returns the path allocated to the given mystery text for saving the result.
    '''
    return 'hdfs://gpu1:8020/results/%s' % mystery_text


def find_result(sc, mystery_text, books_path='hdfs://gpu1:8020/books', min_partitions=8, preserves_partitioning=True, num_partitions=None):
    '''
    Find the mystery text from the books.
    '''
    # books: RDD = [ (book_name_a, book_content_a), (book_name_b, book_content_b), ... ]
    # Returns error if use_unicode is set to False
    books = sc.wholeTextFiles(books_path, minPartitions=min_partitions, use_unicode=True)
    # list_of_book_name_and_sentence: PipelinedRDD = [ (book_name_a, sentence_a_1), (book_name_a, sentence_a_2), (...), (book_name_b, ...), ... ]
    list_of_book_name_and_sentence = books.flatMap(split_content_into_sentences, preservesPartitioning=preserves_partitioning)
    filtered_list_of_book_name_and_sentence = list_of_book_name_and_sentence.filter(lambda x: mystery_text in initialism(x[1]))
    result0 = filtered_list_of_book_name_and_sentence.groupByKey(numPartitions=num_partitions)
    result1 = result0.mapValues(list)
    result2 = result1.collect()
    print('=== Begin Result ===============================================================')
    print(result2)
    print('=== End Result =================================================================\n')
    return result1


def save_result(result, mystery_text):
    '''
    Saves the result.
    '''
    path = get_result_path(mystery_text)
    try:
        # Save the result
        result.saveAsTextFile(path)
    except:
        # Remove the old result file
        remove_existing_file(path)
        result.saveAsTextFile(path)


def load_result(sc, mystery_text):
    '''
    Loads the result.
    '''
    try:
        path = get_result_path(mystery_text)
        result = sc.textFile(path).collect()
        print('=== Begin Result ===============================================================')
        print(result)
        print('=== End Result =================================================================\n')
    except:
        print('*** ERROR: Result does not exist in database ***')


if __name__ == '__main__':
    # mode 0 = generate the result from scratch
    # mode 1 = generate the result from scratch + save result
    # mode 2 = load result from saved result
    mode = int(sys.argv[1])
    if mode == 0 or mode == 1:
        books_path = sys.argv[2]
        mystery_text = str(sys.argv[3]).lower()
        min_partitions = int(sys.argv[4])
        preserves_partitioning = True
        if int(sys.argv[5]) == 0:
            preserves_partitioning = False
        num_partitions = int(sys.argv[6])
        if num_partitions == 0:
            num_partitions = None 
        is_find = True
        if mode == 0:
            is_save = False
        else:
            is_save = True
    elif mode == 2:
        mystery_text = str(sys.argv[2]).lower()
        is_find = False
    else:
        print('*** ERROR: Invalid mode ***')
    # Begin to do real stuff
    print('=== Begin Mystery ==============================================================')
    conf = SparkConf().setAppName('UnravelMysteries').setMaster('yarn')
    sc = SparkContext(conf=conf)
    sc.setLogLevel('INFO')
    if is_find:
        result = find_result(sc, mystery_text, books_path, min_partitions, preserves_partitioning, num_partitions)
        if is_save:
            save_result(result, mystery_text)
    else:
        load_result(sc, mystery_text)
    print('=== End Mystery ================================================================\n')
