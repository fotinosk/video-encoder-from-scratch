## Video compression from scratch

Compression is an old and yet still actively researched topic. In this repo I am replicating some of the older video compression algorithms out there, but trying to optimize them using using some features of Golang.

Modern encoders have achieved upwards of 99% compression ratio, the older methods we are using here can only get us to the ~90% mark. 

To run this:
```bash
cat test_videos/video.rgb24 | go run main.go
```

### The .rgb24 format
A not insignificant challenge of this project, was actually finding an uncompressed video to work with. RGB24 files are structured as matrices containing the RGB values of each pixel for each frame.

### Conversion to YUV
The YUV format is composed of Luminosity (Y) and two chroma values (UV). The underlying assumption of this schema is that the human eye is much more sensitive to small variation in brightness (ie luminosity) than color. Therefore, whilst we store a Y value for each pixel, we now store the color for every 2x2 grid of pixels. 

Conversion from RGB to YUV is a straightforward linear equation. We first calculate the YUV values per pixel and then average the UV values per 2x2 grid of pixels. 

This provides a compression ratio of 50%, which agrees with the theory - RGB format for a 2x2 grid: 2x2x3=12 bytes. For YUV in a 2x2 grid: 4bytes for Y and 2 bytes  for UV (1 for U and one for V) = 6 bytes. Therefore exactly 50% compression. 

### Run Length Encoding

The next important piece of intuition is that frames do not change too much frame-by-frame. This means that adjacent frames will have small difference between them. And hence we can greatly reduce the information we need to store by only storing the difference between adjacent frames. We also need to store keyframes, that act as points of reference. In this case we are only storing the first frame as keynote, from which all differences will be sequentially calculated (frame 4 will be the difference of frame 2 and 3 added on frame 1 - the keyframe). Keeping more keyframes will make the encoding less error-prone, but this is beyond this experiment. 

We start by subtracting every frame from it's previous to only get the temporal difference between frames, and of course storing the first frame as is. We now notice that the calculated difference between frames, is spares, ie mostly made up of zeros. Storing this as is would offer no benefit, as a zero on memory is stored much like any other number, ie taking the same amount of space. 

What we can do, is compress repeatable sets of numbers down using Run-Length Encoding. Here we can store a series of identical numbers using 2 values, first the amount of types the number repeats and then the number itself. For example:

```
[0, 0, 0, 0, 1, 1, 1, 2, 2, 0, 0, 0] => [4, 0, 3, 1, 2, 2, 3, 0]
```
The more repeated numbers a sequence has, the higher the compression ratio. 

This method allows us to compress the video to 14.5% of its original size. 

### Deflate Compression

To see how our algorithm compares with an old, yet well respected compression algorithm, we now try the Deflate Compression Algorithm, to compare results. Deflate is part of most standard libraries, including Go. It should be noted that where the YUV and RLE algorithms we implement are deterministic, giving a consistent value across executions, the Deflate method is not. 

On average, the compression given by Deflate on the selected video is 10% of its original size - pretty close to our algorithm! 

### Time considerations

Large effort was made to optimize the YUV and LRE algorithm, with the LRE compression taking ~50ms, whereas the Deflate algorithm, which is part of the standard library took 13s to compute. That is a ~26x time difference. 
