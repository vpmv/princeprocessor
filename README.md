Prince Processor
---

Resurrected password dictionary generator following [Hashcat Prince Processor](https://github.com/hashcat/princeprocessor).
Corrects the original's overflow errors and outputs a complete and unique dictionary from a custom wordlist.

Contrary to the original, this version generates a list without word-repetition, so passphrases will contain each word only once. <br> 
Next to that output is generated with a word separator (space by default). This aides in applying hashcat rules during a dictionary attack. 

# Install & options

```shell
go install github.com/vpmv/princeprocessor
```

```shell
princeprocessor help

Usage of princeprocessor:
  -m, --max int            Maximum number of elements per chain (default 4)
  -n, --min int            Minimum number of elements per chain (default 2)
  -o, --output string      Output file. Use stdout when omitted
  -s, --separator string   Separator used between elements (default " ")
  -i, --wordlist string    Path to input wordlist file. Use stdin when omitted
```

# How to use

1) Compile a wordlist of possible word combinations, for example:
```text
foo
bar
baz
hello
world
```
2) Run the program

```shell
princeprocessor --min=2 --max=4 -o path/to/output.dict -i path/to/wordlist.txt
```


The program will produce the following output:
```text
foo bar
foo bar baz
foo bar hello
foo bar world
foo baz
foo baz bar
(etc.)
```

> NOTE: For short dictionaries, you can also omit the input file. The program will ask you to input words. 

# Things to consider
The program tries to run as efficiently as possible, using multithreaded parallels. This will utilize relatively more CPU and little RAM. Due to the exponential growth of possible combinations, the program may run for a while. 
In my testing, I've used a wordlist of 60 unique words, with min=2/max=5. Using the formula below, this computes  to (N=2, M=5, n=60) 667,293,540 unique chains.

The output file is around **17Gb**. Because of the exponent M=5, generation takes up to 2 minutes on an old computer. 

$$
S(N, M, n) = \sum_{k=N}^{M} \frac{n!}{(n-k)!}
$$

Numerical equation: `60×59 + 60×59×58 + 60×59×58×57 + 60×59×58×57×56`

In contrast, computing less combinations will be much faster. For example (M=4): 11,912,100 chains; a whopping 655 million less than M=5. Logically, using a higher maximum value will take much longer and be taxing your disk space and I/O. 


## Optional execution (Linux/Mac only)

It may be advantageous to write to RAM instead of disk. This negates disk I/O, making the program much faster, but is **only viable when you have lots of spare RAM**.

For example:

```shell
mkdir /mnt/ramdisk
mount -t tmpfs -o size=16G tmpfs /mnt/ramdisk
```

then run:
```shell
princeprocessor --min=2 --max=5 --wordlist=wordlist.txt -o /mnt/ramdisk/output.dict
```


Found a bug? Please open an issue.
