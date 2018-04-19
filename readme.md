# ProgImage (cloudinary clone)

## Install dependencies
```bash
dep ensure
```

## Run tests
Basic

```bash
go test -v ./...

```

Some tests require a minio server (s3 api), create one easily via docker with
```bash
docker run -p 9000:9000 -e MINIO_ACCESS_KEY=minio -e MINIO_SECRET_KEY=miniostorage minio/minio server /data
```
then set env vars

* `S3_ENDPOINT`
* `S3_ACCESS_KEY`
* `S3_SECRET_KEY`
* `S3_SECURE`

See comments at the top of [s3/image_service_test.go](https://github.com/j0hnsmith/progimage/blob/master/s3/image_service_test.go#L1-L11) for more info.

## Run server
```bash
cd cmd/progimage
go install
progimage server --help
progimage server -a :9090 -e {docker ip}:9000 # or any s3 compatible api
```
See `test.http` for example requests.

