# Unravel the Mysteries by Spark

## Project Description

Mysteries and supenses are the key elements of a crime fiction or movie.  Inside the story, the person behind the crime scene always intentionally leave some clueless-clues there for keep playing the psychological games with the investigator in-charge. **Mystery text** is always the only clue that may uncover the motives, and may allow the investigator to guess what is the next action details of the case. Sometimes, the communications between the gangs will be using mystery text over SMS or email to avoid from traces. Here is the sample of Mystery Text:

`FSASYAOFBFOTCANNCILADTTPTAMACE`
`NWAEIAGCWTWTNOANSCASDCLE`

In real life crime investigation, we understand that as long as we know more about the background and details of such mysterious text, like where it was cited from, the author background, the context of the original story or the original message, etc… we may be able to understand the motives. We may then able to reconstruct the crime scenes and further narrow down the investigation areas.

With this project implementation, the investigator could make use of the program developed to quickly identify whether the "mystery text" seems to be an English intialism or not. If it does cite from somewhere, track the possible sources of the "Quote". The investigator could make use of these information for finding out the truth.

Source: https://blog.cloudera.com/blog/2016/09/solving-real-life-mysteries-with-big-data-and-apache-spark/

## Pros

- Input the mystery text and return the possible sources of the quote in mins;
- Different data sources in different langauages could be incorporated, such as Wikipedia, Project Gutenberg public-domain eBooks, and Bible in order to improve the chances of mystery text resolution. 

## Cons

- As an initial target, the program will resolve the Mystery Text that only constructed in English using Initialism. Definitely more Mystery Text construction methods and data sources in various languages could be incorporated.

## Software to be Installed

- Apache Hadoop
- Apache Spark
- Spark-SQL
- Web Server for Serving the Front-end UI

## Datasets to be Used

- English Wikipedia (JSON dump, https://dumps.wikimedia.org/other/cirrussearch/current/)
- Project Gutenberg public-domain English texts (ftp://mirrors.pglaf.org/mirrors/gutenberg-iso)
- Bible in Basic English (https://www.o-bible.com/dlb.html) 









 







