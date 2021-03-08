
# generate-thumbnails
 generate thumbnails images from mp4, mov.
 and generate vtt file for previews.

 
# Install
 ```shell
 $ git clone https://github.com/yochidros/generate-thumbnails.git | cd generate-thumbnails
 $ ./generate-thumbnails
 ```
or 
```shell
 Download binary from github `generate-thumbnails`
 can use as `$ ./generate-thumbnails`
```

# Usage
Main usage is
```
Usage:
  generate-thumbnails [flags]
    generate-thumbnails [command]

  Available Commands:
    generate    Generate thumbnails
    help        Help about any command

  Flags:
    -h, --help   help for generate-thumbnails
```

Generate usage is
```
Flags:
  -h, --help            help for generate
  -i, --input string    input file path
  -o, --output string   output file path
  -s, --sprit int       thumbnails sprit col length (default 10)
  -t, --time-span int   time span (default 1)
  -w, --width float32   thumbnails width (default 120)
```

# Developement
 - go 1.13.9

# Dependencies
 - FFmpeg 4.3.1 
 - FFprobe 4.3.1 
