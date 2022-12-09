package openai

import (
	"testing"
)

func TestGetGPTTextReply(t *testing.T) {
	reply, err := GetGPTTextReply("你去过的最美的地方是哪")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", reply)
}
