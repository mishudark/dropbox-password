# dropbox-password
Securely stores your passwords as Dropbox do

##Get the library

`go get github.com/mishudark/dropbox-password`

##Encrypt a password
The function needs two arguments, `plaintext password` and the `masterkey`
```
  package main
  
  import "github.com/mishudark/dropbox-password"
  
  func main() {
    hash, err := password.Hash("mishudark", "AES256Key-32Characters1234567890")
  }
```

##Check if a password is valid
The function needs three arguments , `plaintext password`, `hashed password` and `masterkey`

```
  package main
  
  import "github.com/mishudark/dropbox-password"
  
  func main() {
    ok := password.IsValid("mishudark", "aes256$mh68GJ7t9mLYiJKk$7ab2234944dabe98d...", "AES256Key-32Characters1234567890")
  }
```

###Details of implementation

![Image of dropbox]
(https://dropboxtechblog.files.wordpress.com/2016/09/layers.png?w=650&h=443)


It’s universally acknowledged that it’s a bad idea to store plain-text passwords. If a database containing plain-text passwords is compromised, user accounts are in immediate danger. For this reason, as early as 1976, the industry standardized on storing passwords using secure, one-way hashing mechanisms (starting with Unix Crypt). Unfortunately, while this prevents the direct reading of passwords in case of a compromise, all hashing mechanisms necessarily allow attackers to brute force the hash offline, by going through lists of possible passwords, hashing them, and comparing the result. In this context, secure hashing functions like SHA have a critical flaw for password hashing: they are designed to be fast. A modern commodity CPU can generate millions of SHA256 hashes per second. Specialized GPU clusters allow for calculating hashes at a rate of billions per second.

Over the years, we’ve quietly upgraded our password hashing approach multiple times in an ongoing effort to stay ahead of the bad guys. In this post, we want to share more details of our current password storage mechanism and our reasoning behind it. Our password storage scheme relies on three different layers of cryptographic protections, as the figure below illustrates. For ease of elucidation, in the figure and below we omit any mention of binary encoding (base64).

layersMultiple layers of protection for passwords

We rely on bcrypt as our core hashing algorithm with a per-user salt and an encryption key (or global pepper), stored separately. Our approach differs from basic bcrypt in a few significant ways.

First, the plaintext password is transformed into a hash value using SHA512. This addresses two particular issues with bcrypt. Some implementations of bcrypt truncate the input to 72 bytes, which reduces the entropy of the passwords. Other implementations don’t truncate the input and are therefore vulnerable to DoS attacks because they allow the input of arbitrarily long passwords. By applying SHA, we can quickly convert really long passwords into a fixed length 512 bit value, solving both problems.

Next, this SHA512 hash is hashed again using bcrypt with a cost of 10, and a unique, per-user salt. Unlike cryptographic hash functions like SHA, bcrypt is designed to be slow and hard to speed up via custom hardware and GPUs. A work factor of 10 translates into roughly 100ms for all these steps on our servers.

Finally, the resulting bcrypt hash is encrypted with AES256 using a secret key (common to all hashes) that we refer to as a pepper. The pepper is a defense in depth measure. The pepper value is stored separately in a manner that makes it difficult to discover by an attacker (i.e. not in a database table). As a result, if only the password storage is compromised, the password hashes are encrypted and of no use to an attacker.


Source: https://blogs.dropbox.com/tech/2016/09/how-dropbox-securely-stores-your-passwords/
