# crypt

crypt is a small app to encrypt and decrypt files with aes encryption.  It is 
portable and works on OSX, Linux, and Windows along with any other OS that 
supports go.

A key can be specified with the `--key` flag or via stdin if `--key` is 
not specified.  You may want to provide your key via stdin for encryption
and decryption so that it is not saved in your bash history.


```bash
NAME:
   crypt - encrypt and decrypt files easily

USAGE:
   crypt [global options] command [command options] [arguments...]

VERSION:
   1

AUTHOR:
  @crosbymichael

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key                key to use for the encryption algo
   --encrypt, -e        encrypt a file
   --decrypt, -d        decrypt a file
   --help, -h           show help
   --version, -v        print the version
```


## License - MIT
