secureBox
=========
**This is a work in progress.**

A cross platform, end-to-end encrypted file synchronization client in go.

secureBox depends on the following packages: [fsnotify](https://github.com/howeyc/fsnotify).

The goal of secureBox is have the convenience of sharing files between devices and friends without compromising security.


secureBox Encryption spec
=========================
**NOTE:** This document is a work in progress and may not be representative of the final encryption architecture.


Goals
-----
The primary goal of secureBox over other file synchronization services is to have distributed copies of a user's data while making it impossible for a server host or attacker to gain access to said files. This requires files must be unreadable by the remote server or storage host.

Second to this is convenience and the ability to easily share files with other users of the service, ideally regardless of the remote server used by either user. This is possible with the use of symmetric key encryption.


Architecture overview
---------------------

*An understanding of symmetric key encryption is useful to fully understand our encryption methodology.*

Term         | Description
-------------|-------------
File key(s)  | Unique AES key to a particular file.
Key Pair     | RSA Public/private keypair generated at user creation time to ensure the security of a user's File keys.
Password key | AES key dervied from a users password to protect the users master RSA private key.

As stated above, secureBox is designed such that the remote server where a user's files are stored is totally oblivious to said files' content. Passwords, password keys and file keys never leave the users' devices and are never transferred anywhere or to anyone.

The starting point for every decryption process is the user's private key as this one is required to unlock all file file keys. The private key itself however is already encrypted with the user's password which itself never leaves the user's device.


Sharing
-------

Given two users, Alice and Bill, it is easy to securely share a file between the two of them.

* Alice requests Bob's public key from the secureBox Key Server.
* Alice encrypts the file key with Bob's public key.
* Alice writes the new encrypted file key to the encrypted file.
* The secureBox syncs the modified encrypted file.
* Bob uses his private key to decrypt the file key.
* Bob uses the file key to decrypt the file.
