from pyspark import SparkContext
from pyspark.streaming import StreamingContext
import nltk
import sys

def initialism(sentence):
    words = sentence.split()
    initialism = []
    for word in words:
        initialism.append(word[0])
    return ''.join(initialism).lower()

def get_cleaned_sentence_from_zip_txt(book):
    sentences = nltk.tokenize.sent_tokenize(book.replace('\r\n', " ").replace('\n', " "))
    list_match = []
    for sentence in  sentences:
        if user_input.lower() in initialism(sentence):
            list_match.append(sentence.encode('ascii', 'ignore'))
    return list_match
# def check_secret(sentence):


# Create a local StreamingContext with two working thread and batch interval of 10 second
user_input = str(sys.argv[1]).lower()
sc = SparkContext("local[2]", "MysterySearch")
ssc = StreamingContext(sc, 10)
books = ssc.textFileStream("/home/hadoop/data/tmp")
sentences = books.map(get_cleaned_sentence_from_zip_txt).filter(lambda x: x is not None)
sentences.pprint()
ssc.start()
ssc.awaitTermination()