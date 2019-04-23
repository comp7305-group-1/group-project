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


def check_each_book(user_input_lower, book_name_and_list_of_sentences):
    list_of_match = []
    book_name = book_name_and_list_of_sentences[0]
    sentences = book_name_and_list_of_sentences[1]
    for sentence in sentences:  # sentences = sentences list in a book
        if user_input_lower in initialism(sentence):
            list_of_match.append(sentence)
    if len(list_of_match) > 0:
        return {'book_name': book_name, 'sentences': list_of_match}


def get_result_path(user_input):
    return 'hdfs://gpu1:8020/results/%s' % user_input


def find_result0(sc, user_input, min_partitions):
    user_input_lower = user_input.lower()
    # 1. Get sentences from each book
    books = sc.wholeTextFiles('hdfs://gpu1:8020/books', minPartitions=min_partitions, use_unicode=True) # books: RDD; Error after setting use_unicode=False
    book_name_and_list_of_sentences = books.map(split_content_into_sentences)  # book_name_and_list_of_sentences: PipelinedRDD = [ (book_name, [sentence1, sentence2, ...]), (...), (...) ]
    # 2. Get match result from each book
    result0 = book_name_and_list_of_sentences.map(lambda x: check_each_book(user_input_lower, x))
    result = result0.filter(lambda x: x is not None)
    collected = result.collect()
    print('=== Start Result ===============================================================')
    print(collected)
    print('=== End Result =================================================================\n')
    # 3. Save the result to hdfs for later usage
    #path = get_result_path(user_input)
    #try:
    #    # Save the result
    #    result.saveAsTextFile(path)
    #except:
    #    # Remove the old result file
    #    remove_existing_file(path)
    #    result.saveAsTextFile(path)


def find_result(sc, user_input, min_partitions):
    # 1. Get sentences from each book
    books = sc.wholeTextFiles('hdfs://gpu1:8020/testbooks', minPartitions=min_partitions, use_unicode=True) # books: RDD; Error after setting use_unicode=False
    #book_name_and_list_of_sentences = books.map(split_content_into_sentences)  # book_name_and_list_of_sentences: PipelinedRDD = [ (book_name_1, [sentence1, sentence2, ...]), (book_name_2, [...]), ... ]
    # 2. Get match result from each book
    #list_of_book_name_and_sentence = book_name_and_list_of_sentences.flatMapValues(lambda x: x) # list_of_book_name_and_sentence: PipelinedRDD = [ (book_name_1, sentence1), (book_name_1, sentence2), (...), (book_name_2, ...), ...]
    #list_of_book_name_and_sentence_and_initialism = list_of_book_name_and_sentence.map(lambda x: (x[0], x[1], initialism(x[1]))) 
    list_of_book_name_and_sentence = books.flatMap(split_content_into_sentences) # list_of_book_name_and_sentence: PipelinedRDD = [ (book_name_1, sentence1), (book_name_1, sentence2), (...), (book_name_2, ...), ...]
    #result0 = book_name_and_list_of_sentences.map(lambda x: check_each_book(user_input_lower, x))
    user_input_lower = user_input.lower()
    result = list_of_book_name_and_sentence.filter(lambda x: initialism(x[1]) == user_input_lower)
    collected = result.groupByKey().mapValues(list).collect()
    print('=== Start Result ===============================================================')
    print(collected)
    print('=== End Result =================================================================\n')
    # 3. Save the result to hdfs for later usage
    #path = get_result_path(user_input)
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
