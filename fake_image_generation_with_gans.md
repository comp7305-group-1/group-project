# Fake image generation with Generative Adversarial Networks (GANs) with tensorflow on Spark

## Project Description

GANs are neural networks that learn to create synthetic data similar to some known input data. For instance, researchers have generated convincing images from photographs of everything from bedrooms to album covers, and they display a remarkable ability to reflect higher-order semantic logic.

In the project, fake image, for instance, human face, human writing or product,will be generated via Generative Adversarial Networks.

## GAN architecture
Generative adversarial networks consist of two models: a generative model and a discriminative model.
![alt text](https://d3ansictanv2wj.cloudfront.net/GAN_Overall-7319eab235d83fe971fb769f62cbb15d.png "")

The discriminator model is a classifier that determines whether a given image looks like a real image from the dataset or like an artificially created image. This is basically a binary classifier that will take the form of a normal convolutional neural network (CNN).

The generator model takes random input values and transforms them into images through a deconvolutional neural network.

Over the course of many training iterations, the weights and biases in the discriminator and the generator are trained through backpropagation. The discriminator learns to tell "real" images of handwritten digits apart from "fake" images created by the generator. At the same time, the generator uses feedback from the discriminator to learn how to produce convincing images that the discriminator can't distinguish from real images.

## Expected fake image generated:

Human face
![alt text](https://cdn-images-1.medium.com/max/1200/1*-WZvqhgqi1jnpp7wlqjniA.png "")    ![alt text](https://cdn-images-1.medium.com/max/1200/1*7ped99JzASx55f5dIvdkvg.png "")   ![alt text](https://cdn-images-1.medium.com/max/1200/1*FNVgUtYrNZTa8mBtzj7oMw.png "") ![alt text](https://cdn-images-1.medium.com/max/1200/1*xZEOBh8nLHVTDns6ApiASw.png
 "")
 
Human handwriting
![alt text](https://d3ansictanv2wj.cloudfront.net/gan-images-final-d7bdb862726f6fd928a7c859a69c3248.gif "")

Reference: 
https://medium.com/coinmonks/celebrity-face-generation-using-gans-tensorflow-implementation-eaa2001eef86
https://www.oreilly.com/learning/generative-adversarial-networks-for-beginners
## Pros

- There are different GANs resource for case study

## Cons

- Need to research on GAN implementation
- Need to research on GAN and Spark integration 

## Possible Tools

- Apache Hadoop
- Apache Spark
- TensorflowOnSpark
- BigDL
- Web Server

## Datasets

- CelebA dataset (JSON dump, https://dumps.wikimedia.org/other/cirrussearch/current/)
- Stanford Dogs Dataset
 (http://vision.stanford.edu/aditya86/ImageNetDogs/)
- The MNIST Datasbase of handwritten digits 
  (http://yann.lecun.com/exdb/mnist/) 









 






