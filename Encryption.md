secureBox Encryption spec
=========================
**NOTE:** This document is a work in progress and may not be representative of the final encryption architecture.


Goals
-----
The primary goal of secureBox over other file synchronization services is to have distributed copies of a user's data while making it impossible for a server host or attacker to gain access to said files. This requires files must be unreadable by the remote server or storage host. 

Second to this is convenience and the ability to easily share files with other users of the service, ideally regardless of the remote server used by either user. This is possible with the use of symmetric key encryption.


Architecture overview
---------------------
To accomplish the above goals, secureBox will require that every file synchronized will have its own unique public/private key pair. Additonally each user has a master key pair. This allows users a few things:


Given two users, Joe and Bill, it is easy to securely share a file between the two of them.
* Joe gets a copy of Bill's master public key (giving him the ability to encrypt things such that only Bill can see them)
* Joe uses this to encrypt the private key of file foo.txt (and foo.txt itself)
* Joe gives this encrypted version of the private and file key to Bill
* Bill now can both encrypt and decrypt foo.txt ostensibly giving him read/write access


Nuts and Bolts
--------------
*Work in progress.*
