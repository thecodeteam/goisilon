package goisilon

var client *Client
var err error

func init() {
	testClient()
}

func testClient() error {
	client, err = NewClient()
	if err != nil {
		panic(err)
	}
	return nil
}
