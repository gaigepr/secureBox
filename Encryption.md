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
Key Pair     | RSA keypair generated at user creation time to ensure the security of a user's File keys.
Password key | AES key dervied from a users password to protect the users master RSA private key.

Given two users, Alice and Bill, it is easy to securely share a file between the two of them.

* Alice requests Bob's public key from the secureBox Key Server.
* Alice encrypts the file key with Bob's public key.
* Alice writes the new encrypted file key to the encrypted file.
* The secureBox syncs the modified encrypted file.
* Bob uses his private key to decrypt the file key.
* Bob uses the file key to decrypt the file.
