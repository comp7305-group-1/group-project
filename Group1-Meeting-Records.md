**<u>Meeting Minutes</u>**

**1-Mar-2019**

1. All members agreed to research about the term project topics and propose 1-2 topics on 8-Mar.

2. The ideas must meet the following requirements that have been discussed and agreed:
   a. It shouldn't be too complex;

   b. It must allow real-input during demo and should be able to finish within 5 minutes for demonstration purpose; and

   c. Preferrably something new rather than repeating previous projects.

**9-Mar-2019**

1. All members have proposed different topics and try to understand each of them for finalizing the topics;
2. Kenji proposed a topic about generating Rainbow Table for password cracking;
3. Ewen proposed a topic about resolving "mystery text" in crime scene, which the idea came from the website: https://blog.cloudera.com/blog/2016/09/solving-real-life-mysteries-with-big-data-and-apache-spark/
4. Nelson proposed a topic about generating fake images with GANs and TensorFlow;

**15-Mar-2019**

1. Most members agreed to choose "Unravel the Mystery Text" as the term project topic for more interesting and higher confidence levels on implementation.
2. Prof Wang provided a feedback that the topic is interesting, so the term project topic was finalized;
3. Each member will commit the R&R for the term project later.

**18-Mar-2019**

1. Ewen obtained the 3 datasets (8GB ISO File containing open-domain eBooks, 28GB Wikipedia JSON dump, and 4MB Bible English text), and uploaded to the cluster node (student1);
2. All team members discussed the logical flow of the proposed system;
3. Pauline started to look at the data and plan for the cleaning. She mentioned the idea of the cleansing and other members will wait for the result for review;
4. Kenji proposed few more R&R so that he would like to see any team members will be interested in, e.g. Performance Analyst, Report Writer, Coder, etc... 

**20-Mar-2019**

1. Kenji and Ewen started to discuss and plan for the 1st basic configuration for the 8-node cluster;
2. Pauline is still working on the data cleansing; Other team members reviewed and provided input on coming up a smaller, usable dataset.
3. Ewen started to re-configure and standardize the DOM0 and VMs for all student, student2 and student3, such as swap space, memory, vcpus, joining the ganglia, etc...
4. Nelson started to setup an environment for writing the code from his own resources.

**23-Mar-2019**

1. Kenji finalized the settings on GPU1 and reviewed/refined the settings on all 8 VM nodes;
2. Kenji setup user permissions for 'student' and 'hduser' accounts so that that would minimize the chances of mis-configurations or unintentional changes during the project implementation;

**2-Apr-2019**

1. Pauline finished the data cleansing, and explained the result of cleansing to the team. 8GB ISO files were cleaned so that only 3GB texts were left.
2. Kenji and Ewen assisted Pauline to organize the data within HDFS. 
3. Nelson started to test and modify the code on top of the cluster with the cleaned data as the input.

**02-Apr-2019 to 18-Apr-2019 [Ad-hoc Discussions over Whatspp Group with the Following Major Areas]** 

1. Kenji started to study the WebUI integration for demonstration;
2. Ewen worked on the presentation slide decks, with team members' inputs, according to the presentation specification;
3. Nelson came up the version that works on the cluster. Kenji and Ewen performed testing and identify any improve areas to the data processing flow;
4. Kenji and Nelson worked to gather to modify the codes to enhance the parallelism;
5. Ewen performed benchmarking on a freezed version of code with different cluster settings to come up the better performance of the system; e.g. Tried with different Xen configuration & Yarn settings, using Hadoop Archive format (HAR), different HDFS block size to see if performance improved.











