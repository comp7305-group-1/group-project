from pyspark import SparkConf,SparkContext
import re,os,io,sys
import pyhdfs
import zipfile


def func(dirName):
    print(dirName)
    client = pyhdfs.HdfsClient(hosts=hadoopMasterIP,user_name=hadoopMasterName)
    req = client.open(dirName)
    in_memory_data = io.BytesIO(req.data)
    f = zipfile.ZipFile(in_memory_data,'r')
    if len(f.namelist()) == 1:
        for z in f.namelist():
            if '.txt' in z.lower():
                textItem = f.open(z).read()
                zlower = z.lower()
                client.create('/texts/' + zlower,textItem,overwrite=True)
    else:
        print('more than one files in zip')


if __name__ == '__main__':
    hadoopMasterIP = sys.argv[1]
    hadoopMasterName = sys.argv[2]
    sparkMasterURL = sys.argv[3]
     
    # Set up Spark cluster
    conf = SparkConf().setMaster(sparkMasterURL).setAppName('PythonTest')
    sc = SparkContext(conf=conf)
    # Set up Hadoop cluster
    client = pyhdfs.HdfsClient(hosts=hadoopMasterIP,user_name=hadoopMasterName)
    # Enumerate all the folders
    targetZipPath = []
    for i in client.walk('/data'):
        # relative path
        items = i[2]
        print(items)
        # absolute path
        itemsName = [os.path.join(i[0],j) for j in items]
        for j in itemsName:
            if '.zip' in j:
                targetZipPath.append(j)
    # Parallelize these folders through the Spark 
    targetZipPath_ = ['hdfs://' + i for i in targetZipPath]
    targetZipPathRDD = sc.parallelize(targetZipPath_[:300])
    targetZipPathRDD.foreach(func)
