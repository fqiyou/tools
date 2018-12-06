package message

import (
	"github.com/fqiyou/tools/foo/tools/logs"
	"github.com/fqiyou/tools/foo/util"
	"testing"
)


func TestToMessage(t *testing.T){
	log.NewLogger(10000)
	log.SetLevel(log.LevelDebug)
	log.EnableFuncCallDepth(true)

	log.SetLogger("console", "")

	//log.SetLogger("file", fmt.Sprintf(`{"filename":"logs.log","level":%d}`, log.LevelInfo))
	var msg Message

	str := "MDAwMDAwNjd7ImN0aW1lIjoxNTQ0MDExNzYxMDkwLCJpcCI6IjIyMi4yNDQuMjE4LjIxNyIsInByb2plY3QiOiJuZXdzZWFybiJ9ZXlKa2FYTjBhVzVqZEY5cFpDSTZJakUyTmpnMFpHWTNabVpsTWpJdE1EVXpOR0U1TkRSaVlqZ3dPRFF0TlRBd056QmlNemd0TWpNd05EQXdMVEUyTmpnMFpHWTRNREF3WmpjaUxDSnNhV0lpT25zaUpHeHBZaUk2SW1weklpd2lKR3hwWWw5dFpYUm9iMlFpT2lKamIyUmxJaXdpSkd4cFlsOTJaWEp6YVc5dUlqb2lNUzQzTGpFMEluMHNJbkJ5YjNCbGNuUnBaWE1pT25zaUpITmpjbVZsYmw5b1pXbG5hSFFpT2pZME1Dd2lKSE5qY21WbGJsOTNhV1IwYUNJNk16WXdMQ0lrYkdsaUlqb2lhbk1pTENJa2JHbGlYM1psY25OcGIyNGlPaUl4TGpjdU1UUWlMQ0lrYkdGMFpYTjBYM0psWm1WeWNtVnlJam9pSWl3aUpHeGhkR1Z6ZEY5eVpXWmxjbkpsY2w5b2IzTjBJam9pSWl3aVlXUmZjR0ZuWlNJNkltNWxkM05mWkdWMFlXbHNJaXdpWVdSZmNHOXphWFJwYjI0aU9qRXdMQ0poWkY5aFkzUnBiMjRpT2lKbGVIQnZjM1Z5WlNJc0ltRmtYMmxrSWpvaWRUTTFNVFl3TVRVaUxDSmhaRjkwZVhCbElqb2lNaTB4T1NJc0luQmhaMlZmZFhKc0lqb2lhSFIwY0hNNkx5OTNkM2N1ZUdsdWQyVnVlbWgxWVc0dWJtVjBMM05vWVhKbEwzaDNlbDloY25ScFkyeGxMMkZ5ZEdsamJHVmZhSFIwY0hNMFgyeHZibWN1YUhSdGJDSXNJbU5vWVc1dVpXd2lPaUlpTENKemIzVnlZMlVpT2lJMElpd2lZMkYwWldkdmNua2lPakFzSW05eWFXZHBiaUk2SWpRaUxDSjFjMlZ5U1dRaU9pSXlPVGs1T1RNME1pSXNJblZoSWpvaWJXOTZhV3hzWVM4MUxqQWdLR3hwYm5WNE95QmhibVJ5YjJsa0lEVXVNUzR4T3lCeVpXUnRhU0F6SUdKMWFXeGtMMnh0ZVRRM2Rqc2dkM1lwSUdGd2NHeGxkMlZpYTJsMEx6VXpOeTR6TmlBb2EyaDBiV3dzSUd4cGEyVWdaMlZqYTI4cElIWmxjbk5wYjI0dk5DNHdJR05vY205dFpTODFOeTR3TGpJNU9EY3VNVE15SUcxdlltbHNaU0J6WVdaaGNta3ZOVE0zTGpNMklpd2lKRzl6SWpvaVlXNWtjbTlwWkNJc0lpUmtaWFpwWTJWZmFXUWlPaUl5T1RrNU9UTTBNaUlzSWlScGJXVnBJam9pTWprNU9Ua3pORElpTENJa2FYTmZabWx5YzNSZlpHRjVJanBtWVd4elpYMHNJblI1Y0dVaU9pSjBjbUZqYXlJc0ltVjJaVzUwSWpvaVYyVmlRV1JFWVhSaElpd2lYMjV2WTJGamFHVWlPaUl3T1RjeE5qSXhOakkzTVRRMU5qSWlmUT09"
	str = "MDAwMDAwNjd7ImN0aW1lIjoxNTQ0MDExNzYxMDkwLCJpcCI6IjIyMi4yNDQuMjE4LjIxNyIsInByb2plY3QiOiJuZXdzZWFybiJ9ZXlKa2FYTjBhVzVqZEY5cF"

	msg.ToMessageList(str)

	util.JsonPrint(msg.MessageList)
	util.JsonPrint(msg.MessageInfo)


}