# Redisee

This is a tool to analysis the memory usage of the whole redis instance.
Use [github.com/redis/go-redis/v9](https://github.com/redis/go-redis?tab=readme-ov-file) underneath

# Installation
The binary files have been released, you can get it through the following command:
```bash
wget {binary file source link} -O redisee
```

| OS | Arch | Binary File Source Link|
| --- | --- | --- |
| darwin | amd64 | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_darwin_amd64 |
| darwin | arm64 | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_darwin_arm64 |
| linux | 386 | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_linux_386 |
| linux | amd64 | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_linux_amd64 |
| linux | arm | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_linux_arm |
| linux | arm64 | https://github.com/butflame/redisee/releases/download/v1.0.0/redisee_linux_arm64 |

If you are using Windows, you can find the binary file in the release page.

If you are using other OS or Arch, you can build it yourself.

# Usage
The binary files are executable files, you can use it directly.

```bash
./redisee
```

You can check the help message by the following command:
```bash
./redisee -h
```

All available flags are as follows:
| Flag | Default Value | Description |
| --- | --- | --- |
| -h |  | Show help message |
| -host | 127.0.0.1 | the host of the redis server |
| -port | 6379 | the port of the redis server |
| -password |  | password to do authorization |
| -pattern | * | the pattern to use with `scan` command |
| -seq | any non-numeric and non-alphabetic character  | seperators to determine the prefix of key, **every character** of this flag value would be used as seperator. |
| -concurrency | 4 | concurrency when executing `ttl` and `memory usage` for keys. |

Here's a showcase of the output:

At first it will echo the settings.
```bash
➜  ~ ./redisee
Start with setting:
Host: 127.0.0.1
Port: 6379
Password: 
Db: All
Separator: 
Scan pattern: *
Concurrency: 4

you can run with "-h" to see available flags
```

Then it will start to scan by db.
```bash
➜  ~ ./redisee

...

Scanning 6 db(s), 4 concurrency, scan pattern: *
Scanning db 2, keys 142268/142357 scanned
Scanning db 3, keys 140833/140914 scanned
Scanning db 4, keys 143226/143311 scanned
Scanning db 5, keys 188134/188279 scanned
Scanning db 0, keys 134728/134799 scanned
Scanning db 1, keys 141097/141181 scanned
Finished scanning all dbs
```

Finally it will print the result, including:
- Overall memory usage and total keys
- Key count by ttl
- Top100 keys by memory usage
- Total memory usage by key type
- Top100 key prefix by key count
- Top100 key prefix by memory usage
```bash
### Overall:
Total Memory: 0 B, Used Memory: 404.61 MB, Total Keys: 890286

### Total Keys By TTL:
no exp: 2108 keys in total
expired: 0 keys in total
0-1h: 5373 keys in total
1-3h: 10443 keys in total
3-12h: 48513 keys in total
12-24h: 63870 keys in total
1-2day: 128645 keys in total
3-7day: 631334 keys in total
>7day: 0 keys in total

### Top100 Keys by Memory Usage:
Key: ifl:set:apple2:cow2:dog:1045696407, Type: set, Memory Usage: 1.48 KB
Key: ifl:set:apple2:cow2:dog2:1561508988, Type: set, Memory Usage: 1.48 KB
Key: ifl:set:apple2:cow2:dog2:9552488249, Type: set, Memory Usage: 1.48 KB
...

### Statistics By Key Type:
Type: string, Count: 386975, Memory Usage: 33.92 MB
Type: hash, Count: 119147, Memory Usage: 68.06 MB
Type: list, Count: 117950, Memory Usage: 46.18 MB
Type: set, Count: 156024, Memory Usage: 156.69 MB
Type: zset, Count: 110190, Memory Usage: 40.60 MB

### Top100 Key Prefix By Key Count:
Prefix: ifl:, Count: 890286, Memory Usage: 345.45 MB
Prefix: ifl:string:, Count: 386975, Memory Usage: 33.92 MB
Prefix: ifl:string:ab3:, Count: 281253, Memory Usage: 25.21 MB
...

### Top100 Key Prefix By Memory Usage:
Prefix: ifl:, Count: 890286, Memory Usage: 345.45 MB
Prefix: ifl:set:, Count: 156024, Memory Usage: 156.69 MB
Prefix: ifl:set:apple2:, Count: 111416, Memory Usage: 112.13 MB
...

```

# Tips
The binary files are not tested on Windows.

Be careful when using this in a production environment. If the `-concurrency` flag is set too high, it will exhaust the network or the redis server may be overloaded.

# License
MIT
