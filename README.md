# Huggingface CLI in Golang
A simple CLI tool to generate images from text using the Huggingface API in Golang.

## Installation
```bash
go install github.com/gnzdotmx/huggingface-cli/huggingface-cli@latest
```
Troubleshooting:
```
go clean -modcache
export PATH=$PATH:`go env GOPATH`/bin
```
## Configuration
- Create an API key at https://huggingface.co/account
- Set the `huggingface_api_key` environment variable with your Huggingface API key.
```bash
export hugingface_api_key=YOUR_API
```
## Usage
```bash
huggingface-cli -f input.txt -m black-forest-labs/FLUX.1-schnell
```
- `-f` flag specifies the text file to read.
- `-m` flag specifies the model to use.

## Text file example
### input.txt
This tool will generate one image per line 
```
The quick brown fox jumps over the lazy dog.

The sky was clear and the stars were shining brightly.

The old castle stood on the hill, overlooking the village.
```

### Output

All images have been generated and saved.
```
generated_image_1.png
generated_image_2.png
generated_image_3.png
```