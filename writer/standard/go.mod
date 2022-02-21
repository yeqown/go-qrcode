module github.com/yeqown/go-qrcode/writer/standard

go 1.17

require (
	github.com/fogleman/gg v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/yeqown/go-qrcode/v2 v2.2.0
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

//replace github.com/yeqown/go-qrcode/v2 => ../../
