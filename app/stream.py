import sys
# must for python 2
reload(sys)
sys.setdefaultencoding('utf-8')

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
    for sentence in sentences:
        if user_input.lower() in initialism(sentence):
            list_match.append(sentence.encode('ascii', 'ignore'))
    return list_match
""" Setup
# 1. run this program /opt/spark-2.4.0-bin-hadoop2.7/bin/spark-submit --mas
ter=yarn --py-files ../application/dependencies.zip s.py tpgeo
# 2. cp a new txt to /home/hduser/tmp in this case
# 3. New results will be shown on the console
# Remarks, input sentences should not contain \n if they are in 1 sentences, otherwise spark will stream line by line which may cause wrong result.
"""

# Create a local StreamingContext with two working thread and batch interval of 10 second
user_input = str(sys.argv[1]).lower()
# N thread, here is 2
sc = SparkContext("local[2]", "MysterySearch")
# Check folder every 8 second
ssc = StreamingContext(sc, 8)
# Must user "file://" + directory
books = ssc.textFileStream("file:///home/hduser/tmp")
sentences = books.map(get_cleaned_sentence_from_zip_txt).filter(lambda x: x is not None)
sentences.pprint()
# Start monitoring
ssc.start()
# terminate with the program
ssc.awaitTermination()
