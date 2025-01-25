# tiktok-ecommerce

## Folder structure

```base
├── README.md
├── app
│   └── user
├── go.work
├── go.work.sum
└── idl
    └── user.proto
```

### Use [cwgo](https://github.com/cloudwego/cwgo) to create boilerplate code for RPC service

```base
# at the root folder

# install cwgo
go install github.com/cloudwego/cwgo@latest
go env #check out GOPATH, and add it to $PATH to enable cwgo in your terminal

# create idl using protobuf
cd idl
touch user.proto

# create user service boilerplate code
mkdir -p app/user
cd app/user
cwgo server -I ../../idl --type RPC --module github.com/jazwu/tiktok-ecommerce/app/user --service user --idl ../../idl/user.proto

go mod tidy # update dependencies
go work use . # add workspace
go run .
```

