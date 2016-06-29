# carve
Go implementation of Seam Carving algorithm. Seam carving allows for an image to be resized without simply cropping or scaling the original image. Detailed walk through of the algorithm and code can be found on [my website](http://parellagram.com/posts/carve). Algorithm is also detailed on [Wikipedia](https://en.wikipedia.org/wiki/Seam_carving), the [original paper](http://graphics.cs.cmu.edu/courses/15-463/2007_fall/hw/proj2/imret.pdf). 

# Usage

    resized := carve.ReduceHeight(img, 50)

# Todo

* Implement vertical seam carving
* Improve performance by caching?
