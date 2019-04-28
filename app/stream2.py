import sys
# must for python 2
reload(sys)
sys.setdefaultencoding('utf-8')

from pyspark import SparkContext, SparkConf
from pyspark.streaming import StreamingContext
import nltk
import sys


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


# def get_cleaned_sentence_from_zip_txt(book):
#     sentences = nltk.tokenize.sent_tokenize(book.replace('\r\n', " ").replace('\n', " "))
#     list_match = []
#     for sentence in sentences:
#         if user_input.lower() in initialism(sentence):
#             list_match.append(sentence.encode('ascii', 'ignore'))
#     return list_match


""" Setup
# 1. run this program /opt/spark-2.4.0-bin-hadoop2.7/bin/spark-submit --mas
ter=yarn --py-files ../application/dependencies.zip s.py tpgeo
# 2. cp a new txt to /home/hduser/tmp in this case
# 3. New results will be shown on the console
# Remarks, input sentences should not contain \n if they are in 1 sentences, otherwise spark will stream line by line which may cause wrong result.
"""
if __name__ == '__main__':
    books_path = sys.argv[1]
    mystery_text = str(sys.argv[2]).lower()
    batch_duration = int(sys.argv[3])
    preserves_partitioning = True
    if int(sys.argv[4]) == 0:
        preserves_partitioning = False

    # Begin to do real stuff
    print('=== Begin Mystery (Streaming) ==================================================')
    conf = SparkConf().setAppName('UnravelMysteriesStream').setMaster('yarn')
    sc = SparkContext(conf=conf)

    # Check folder every batch_duration seconds
    ssc = StreamingContext(sc, batch_duration)

    # Must user "file://" + directory
    #books = ssc.textFileStream("file:///home/hduser/tmp")
    books = ssc.textFileStream(books_path)

    #sentences = books.map(get_cleaned_sentence_from_zip_txt).filter(lambda x: x is not None)
    list_of_book_name_and_sentence = books.flatMap(split_content_into_sentences, preservesPartitioning=preserves_partitioning)
    filtered_list_of_book_name_and_sentence = list_of_book_name_and_sentence.filter(lambda x: mystery_text in initialism(x[1]))

    #sentences.pprint()
    filtered_list_of_book_name_and_sentence.pprint()

    # Start monitoring
    ssc.start()
    # terminate with the program
    ssc.awaitTermination()
    print('=== End Mystery (Streaming) ====================================================\n')
