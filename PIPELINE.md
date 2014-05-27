File transformation pipeline
============================

The Goal of this document is to capture implementation details of the file transformation pipeline.

Actors in the transformation pipeline will implement PipeReader and/or PipeWriter.
http://golang.org/src/pkg/io/pipe.go

Target actors to implement:
---------------------------

Actor Name | PipeReader | PipeWriter | Purpose
---------- | ---------- | ---------- | -------
BufferedFile | X | X | Read/Write data to/from local filesystem in a buffered manner
FixedBlockChunker | | X | Break streamed bytes into smaller X sized chunks
SlidingWindowChunker | | X | Break streamed bytes into smaller sized chunks using rabin fingerprinting
ContainerFormat | X | X | Control reading/writing data into a structured container at rest
CompressedBlock | X | X | Compress/Decompress data read/write through actor
DedupeBlock | X | X | Control validation of duplicate data and pass through reference markers
EncryptedBlock | X | X | Encrypt/Decrypt data read/write through actor
BufferedHTTP | X | X | Read/Write data to/from HTTP object storage (S3 like)
HashCollector | X | X | Calculate running hashes of data read/write through actor
ReedSolomnBlock | X | X | Calculate data to build reed solomn code blocks of data read/write through actor

As a file enters the system for copy a read pipe and write pipe will be constructed.
This construction will be driven based on metadata captured from the file source and destination.

An example:
-----------

* source file location ~/secure-box/family photos/trip to disney/mickey_is_goofy.jpg
* destination file location /mnt/backup
* metadata collected: local file source, local file destination, option flag to compress, option flag to chunk, option flag to encrypt

```Actors constructed from left to right, right most actor is point of entry to the pipeline```

read pipe:
----------
BufferedFile | HashCollector

write pipe:
-----------
BufferedFile | HashCollector | ContainerFormat | EncryptedBlock | CompressedBlock | HashCollector | FixedBlockChunker

The process to copy will read X amount from the read pipe and write X amount to write pipe.
The bytes will be transformed to format they will become at rest through this pipe.

In reverse to "restore" the file the pipe would like

read pipe:
----------
BufferedFile | HashCollector | ContainerFormat | EncryptedBlock | CompressedBlock | HashCollector

write pipe:
-----------
BufferedFile | HashCollector


Container Format
================

item name | item purpose
--------- | ------------
magic_number | http://lmgtfy.com/?q=magic+numbers
container_version | 4 bits to represent the version of the container used
container_flags | 64 bits to represent metadata for setting up a read pipeline
header_start | guid marker to represent the start of the header section
header_field_size | 4 bits to represent the size of the field to read in
header_field_value | value to be read in
... | lather, rinse, repeat
header_size | 32 bits to represent the overall size of the header container
header_hash | hash of data written to header for validation against data corruption
data_start | guid marker to represent the start of the data section
block_type | 4 bits to represent the type of the block
src_size | 32 bits to present the size of the data read from disk
block_id | 64 bits to represent the unique id of the data
block_size | 32 bits to represent the size of block to read from the container
block_start | begin data section
... | lather, rinse, repeat
footer_start | guid marker to represent the start of the footer section
footer_field_size | 4 bits to represent the size of the field to read in
footer_field_value | value to be read in
... | lather, rinse, repeat
footer_size | 32 bits to represent the overall size of the footer container
footer_hash | hash of data written to footer for validation against data corruption
