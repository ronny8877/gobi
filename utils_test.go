// SPDX-License-Identifier: MIT

package main

import "testing"

func TestParseArguments(t *testing.T) {
	args, err := parseArguments("len=2,type=Array(len=3,type=Int(min=1,max=2)),text=https://example.com/?a=b")
	if err != nil {
		t.Fatal(err)
	}
	if args["len"] != "2" || args["type"] != "Array(len=3,type=Int(min=1,max=2))" || args["text"] != "https://example.com/?a=b" {
		t.Fatal(args)
	}
	for _, input := range []string{"len", "len=2,len=3", "type=Int())", "type=Int(", "=x"} {
		if _, err := parseArguments(input); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
}
func TestGeneratorBounds(t *testing.T) {
	result := responseBuilderWithRefs(map[string]any{"values": "Array(len=10000,type=Array(len=10000,type=Int()))"}, nil, 0)
	if len(result["values"].([]any)) > maxArrayLength {
		t.Fatal("unbounded output")
	}
}
