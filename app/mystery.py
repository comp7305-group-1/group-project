import os
import zipfile
import sys
from nltk.tokenize import sent_tokenize
from pyspark import SparkContext, SparkConf

zip_dir = "/home/hadoop/data/10001/"
# Step 1 Get all zips from a directory "path" , change in different computer
def get_zips():
    files = []
    for r, d, f in os.walk(zip_dir):
        for file in f:
            if '.ZIP' in file:
                files.append(os.path.join(r, file))
    return files

# Step 2 Get a txt from a zip into a cleaned sentence
def get_cleaned_sentence_from_zip_txt(zip_dir="/home/hadoop/data/10001/10001.ZIP"):
    archive = zipfile.ZipFile(zip_dir)
    _, zip_name = os.path.split(zip_dir)
    zip_file = archive.read(zip_name.replace('.ZIP', "").replace('.zip', '') + '.txt')
    return sent_tokenize(zip_file.replace('\r\n', " ").replace('\n', " "))


# # Step 3 Get Sentences from book
# def clean_book_and_get_sentences(fs, book):
#     lines = sc.textFile("/home/hadoop/data/10001/10001.txt")
#     #TODO cleaning definition
#     #Return sentences
#     return book

# Step 4 Get the initials of a sentence
def initialism(sentence):
    words = sentence.split()
    initialism = []
    for word in words:
        initialism.append(word[0])
    return ''.join(initialism).lower()

def start(user_input="tpg", total_count=1):
    count = 0
    print("Start Mystery===============================================")
    # Step 1
    zips_dir = get_zips()
    print("Number of zips found: %s" % len(zips_dir))
    # Step 2
    for zip_dir in zips_dir:
        # Step 3
        sentences = get_cleaned_sentence_from_zip_txt(zip_dir)
        for sentence in sentences:
            # Step 4
            if user_input.lower() in initialism(sentence):
                count += 1
                print("Possible mystery found %s: \n%s" % (count, sentence))
                if count >= total_count:
                    print("End Mystery===============================================")
                    return

if __name__ == '__main__':
    conf = SparkConf().setAppName("app")
    sc = SparkContext(conf=conf)
    start(str(sys.argv[1]), int(sys.argv[2]))
