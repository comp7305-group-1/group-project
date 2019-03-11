# Password Cracking using Rainbow Table on Spark

### Project Description

The objective of this project is to build a cryptographic hash solver based on rainbow table. The user inputs the cryptographic hash, like MD5 or SHA1, the cluster tries to find the original text using a rainbow table.

### Pros

- The problem well suited to be computed in a distributed environment for speed up.
- The dataset size could be adjusted at will, because the rainbow table itself is generated. Several parameters control its size:
  - The target domain to be solved (e.g. numberic strings, alphanumeric string, etc.)
  - The length of a chain
  - The number of chains

### Cons

- It is difficult to distribute work amoung 4 groupmates.
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
