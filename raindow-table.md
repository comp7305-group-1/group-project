# Cryptographic Hash Solver on Spark

### Project Description

Cryptographic hash solver using rainbow table. The user inputs the cryptographic hash, like MD5 or SHA1, the cluster tries to find the original text using a rainbow table.

### Pros

- The problem is suitable for being computed in a distributed environment for speed up.
- The dataset size could be adjusted freely, since the rainbow table itself is generated, and several parameters control its size.

### Cons

- Difficult to distribute tasks to 4 groupmates.
- The performance may not be fascinating on a 4-machine cluster.

### Software to be Installed

- Apache Hadoop ( https://hadoop.apache.org/ ):
  - For storing the rainbow table
- Apache Spark ( https://spark.apache.org/ ):
  - For computing the solution

### Open-Source Programs

(none)

### Dataset(s) to be Used

- Rainbow table generated from this project in "Generation Mode"
