# stagecoach

A server to upload files to a specific directory via HTTP.

## Example Config

The config file should reside in `/etc/stagecoach/stagecoach.yml` or in the current working directory.

`authorizationToken` is a sha512 hash used to authorize requests, it's like a password. (Example is "MyTopSecretToken")

```yaml
destinationDirectory: /tmp/test
authorizationToken: fdd07c6da3405b29aa912e39662fd2cb6e19fde35ec9a81c708fcfb265aff7a15ad0c5cbd87370bf096035bc43a58351c437e6ef599077bcbfbbf0aabcbfd770
```

## Running via docker

Start the container with

```shell
docker run --rm \
  -v $(pwd)/stagecoach.yml:/etc/stagecoach/stagecoach.yml \
  -p 4444:4444 \
  ghcr.io/nitwhiz/stagecoach:latest
```

## Uploading files

After starting the server, you can upload files via `POST http://localhost:4444/upload`:

```shell
curl --request POST \
  --url http://localhost:4444/upload \
  --header 'Authorization: Token MyTopSecretToken' \
  --header 'content-type: multipart/form-data' \
  --form file=@file \
  --form 'name=my_test_file.txt'
```

Existing files are not overwritten.
