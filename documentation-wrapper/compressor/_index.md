+++
title = "compressor"
chapter = false
weight = 5
+++
![logo](/wrappers/compressor/compressor.svg?width=200px)
## Summary

Compressor is a wrapper that takes any binary (any file type or OS) and packages it into the specified archive format. 

### Supported Methods
This wrapper uses the common OS commands to archive the binary.

There are three options when building
- The `name` specifies the file name of the payload inside the compressed archive
- The `method` specifies the method to use to compress the payload
- The `password` options is optional and only used if the `zip` method is selected

```bash
# Normal zip
zip -j package.zip agent.exe

# Zip with a password
zip -jeP password package.zip agent.exe

# tar
tar -cvf package.tar agent.exe

# tar (Gzip)
tar -czvf package.tar.gz agent.exe

# tar (Bzip2)
tar -cjvf package.tar.bz2 agent.exe

# tar (XZ)
tar -cJvf package.tar.xz agent.exe
```

## Author
- @lum8rjack
