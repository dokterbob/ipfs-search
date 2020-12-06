package ipfs

import (
	"context"
	"fmt"
	"github.com/dankinder/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

type StatTestSuite struct {
	suite.Suite

	ctx  context.Context
	ipfs *IPFS

	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
}

func (s *StatTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)

	cfg := DefaultConfig()
	cfg.IPFSAPIURL = s.mockAPIServer.URL()

	s.ipfs = New(cfg, http.DefaultClient, instr.New())
}

func (s *StatTestSuite) TearDownTest() {
	s.mockAPIServer.Close()
}

func (s *StatTestSuite) TestDirectory() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv",
		},
	}

	rURL := fmt.Sprintf("/api/v0/files/stat?arg=%%2Fipfs%%2F%s", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`{"Hash":"QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv","Size":0,"CumulativeSize":6544,"Blocks":7,"Type":"directory"}`),
		}).
		Once()

	err := s.ipfs.Stat(s.ctx, r)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Equal(r.Stat, t.Stat{
		Type: t.DirectoryType,
		Size: 6544,
	})
}

func (s *StatTestSuite) TestSinglePartFile() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB",
		},
	}

	rURL := fmt.Sprintf("/api/v0/files/stat?arg=%%2Fipfs%%2F%s", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`{"Hash":"QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB","Size":1091,"CumulativeSize":1102,"Blocks":0,"Type":"file"}`),
		}).
		Once()

	err := s.ipfs.Stat(s.ctx, r)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Equal(r.Stat, t.Stat{
		Type: t.FileType,
		Size: 1091,
	})
}

func (s *StatTestSuite) TestMultiPartFile() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
		},
	}

	rURL := fmt.Sprintf("/api/v0/files/stat?arg=%%2Fipfs%%2F%s", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`{"Hash":"Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi","Size":4475792,"CumulativeSize":4476917,"Blocks":18,"Type":"file"}`),
		}).
		Once()

	err := s.ipfs.Stat(s.ctx, r)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Equal(r.Stat, t.Stat{
		Type: t.FileType,
		Size: 4475792,
	})
}

func (s *StatTestSuite) TestUnreferencedPartial() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
	}

	rURL := fmt.Sprintf("/api/v0/files/stat?arg=%%2Fipfs%%2F%s", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`{"Hash":"QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p","Size":262144,"CumulativeSize":262158,"Blocks":0,"Type":"file"}`),
		}).
		Once()

	err := s.ipfs.Stat(s.ctx, r)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Equal(r.Stat, t.Stat{
		Type: t.PartialType,
		Size: 262144,
	})
}

func (s *StatTestSuite) TestReferencedPartial() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			},
			Name: "fileName.pdf",
		},
	}

	rURL := fmt.Sprintf("/api/v0/files/stat?arg=%%2Fipfs%%2F%s", r.ID)

	// Setup mock handler
	s.mockAPIHandler.
		On("Handle", "POST", rURL, mock.Anything).
		Return(httpmock.Response{
			Body: []byte(`{"Hash":"QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p","Size":262144,"CumulativeSize":262158,"Blocks":0,"Type":"file"}`),
		}).
		Once()

	err := s.ipfs.Stat(s.ctx, r)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Equal(r.Stat, t.Stat{
		Type: t.FileType,
		Size: 262144,
	})
}

// TODO: TestInvalid()

// TODO: TestGeneralError()

func TestStatTestSuite(t *testing.T) {
	suite.Run(t, new(StatTestSuite))
}

// Protobuf ref: http://docs.ipfs.io.ipns.localhost:8080/concepts/file-systems/#unix-file-system-unixfs
// Multicodec type reference: https://github.com/multiformats/multicodec/blob/master/table.csv

// Mocking HTTP
// https://github.com/dankinder/httpmock
// https://hassansin.github.io/Unit-Testing-http-client-in-Go
// https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/

// Anything not a timeout -> invalid (?)

// Invalid:
// bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq -> file/ls: expected protobuf dag node
// QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8 -> file/ls: proto: required field "Type" not set
// QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj -> file/ls: proto: unixfs_pb.Data: illegal tag 0 (wire type 0)

// Unsupported in Ls (but supported in Stat!):
// QmToQ3m6g8XdnMhoMR2hdxrvFtKAEX2DMcWpnFM6YifXQD -> file/ls: unrecognized type: Raw

// curl -q -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq" | jq
// {
//   "Hash": "bafkreia2whgx2vblgdpwim5ugz7ofhxoo2vtpyart633mj6gbpwsj7yfxq",
//   "Size": 262144,
//   "CumulativeSize": 262144,
//   "Blocks": 0,
//   "Type": "file"
// }

// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8 HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:05:48 GMT
// < Transfer-Encoding: chunked
// <
// { [88 bytes data]
// 100    77    0    77    0     0  59597      0 --:--:-- --:--:-- --:--:-- 77000
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Message": "proto: required field \"Type\" not set",
//   "Code": 0,
//   "Type": "error"
// }

// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:06:15 GMT
// < Transfer-Encoding: chunked
// <
// { [100 bytes data]
// 100    89    0    89    0     0  24211      0 --:--:-- --:--:-- --:--:-- 29666
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Message": "proto: unixfs_pb.Data: illegal tag 0 (wire type 0)",
//   "Code": 0,
//   "Type": "error"
// }

// QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq -> file/ls: unexpected EOF

// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:06:48 GMT
// < Transfer-Encoding: chunked
// <
// { [64 bytes data]
// 100    53    0    53    0     0  46207      0 --:--:-- --:--:-- --:--:-- 53000
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Message": "unexpected EOF",
//   "Code": 0,
//   "Type": "error"
// }

// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/QmToQ3m6g8XdnMhoMR2hdxrvFtKAEX2DMcWpnFM6YifXQD" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/QmToQ3m6g8XdnMhoMR2hdxrvFtKAEX2DMcWpnFM6YifXQD HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 200 OK
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:07:10 GMT
// < Transfer-Encoding: chunked
// <
// { [127 bytes data]
// 100   121    0   121    0     0  25398      0 --:--:-- --:--:-- --:--:-- 30250
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Hash": "QmToQ3m6g8XdnMhoMR2hdxrvFtKAEX2DMcWpnFM6YifXQD",
//   "Size": 262144,
//   "CumulativeSize": 262158,
//   "Blocks": 0,
//   "Type": "file"
// }

// Ethereum block
// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:09:57 GMT
// < Transfer-Encoding: chunked
// <
// { [79 bytes data]
// 100    68    0    68    0     0  55737      0 --:--:-- --:--:-- --:--:-- 68000
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Message": "unrecognized object type: 144",
//   "Code": 0,
//   "Type": "error"
// }

// Git repo
// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 500 Internal Server Error
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:10:29 GMT
// < Transfer-Encoding: chunked
// <
// { [80 bytes data]
// 100    69    0    69    0     0  63129      0 --:--:-- --:--:-- --:--:-- 69000
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Message": "not unixfs node (proto or raw)",
//   "Code": 0,
//   "Type": "error"
// }

// Correct directory
// $ curl -v -X POST "http://127.0.0.1:5001/api/v0/files/stat?arg=/ipfs/QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D" | jq
//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
//                                  Dload  Upload   Total   Spent    Left  Speed
//   0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 5001 (#0)
// > POST /api/v0/files/stat?arg=/ipfs/QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D HTTP/1.1
// > Host: 127.0.0.1:5001
// > User-Agent: curl/7.54.0
// > Accept: */*
// >
// < HTTP/1.1 200 OK
// < Access-Control-Allow-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Access-Control-Expose-Headers: X-Stream-Output, X-Chunked-Output, X-Content-Length
// < Content-Type: application/json
// < Server: go-ipfs/0.6.0
// < Trailer: X-Stream-Error
// < Vary: Origin
// < Date: Sun, 06 Dec 2020 14:11:00 GMT
// < Transfer-Encoding: chunked
// <
// { [132 bytes data]
// 100   126    0   126    0     0   111k      0 --:--:-- --:--:-- --:--:--  123k
// * Connection #0 to host 127.0.0.1 left intact
// {
//   "Hash": "QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D",
//   "Size": 0,
//   "CumulativeSize": 61702258918,
//   "Blocks": 6,
//   "Type": "directory"
// }

// Parent of a multiblock file
// {
//   "Hash": "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
//   "Size": 4475792,
//   "CumulativeSize": 4476917,
//   "Blocks": 18,
//   "Type": "file"
// }

// Second part of a multiblock file
// {
//   "Hash": "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
//   "Size": 262144,
//   "CumulativeSize": 262158,
//   "Blocks": 0,
//   "Type": "file"
// }

// Long directory (wikipedia)
// {
//   "Hash": "QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp",
//   "Size": 0,
//   "CumulativeSize": 613715579624,
//   "Blocks": 256,
//   "Type": "directory"
// }