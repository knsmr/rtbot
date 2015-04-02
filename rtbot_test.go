package main

import "testing"

func TestPageUrl(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{1, "http://jp.techcrunch.com"},
		{2, "http://jp.techcrunch.com/page/2/"},
		{3, "http://jp.techcrunch.com/page/3/"},
		{4, "http://jp.techcrunch.com/page/4/"},
	}
	for _, c := range cases {
		got := pageUrl(c.in)
		if got != c.want {
			t.Errorf("pageUrl(%d) == %q, want %q", c.in, got, c.want)
		}
	}
}

func Test_TweetWorthy(t *testing.T) {
	cases := []struct {
		in   [2]int
		want bool
	}{
		{[2]int{0, 0}, false},
		{[2]int{10, 3}, false},
		{[2]int{50, 0}, true},
		{[2]int{50, 50}, false},
		{[2]int{99, 0}, true},
		{[2]int{100, 0}, true},
		{[2]int{100, 100}, false},
		{[2]int{103, 98}, true},
		{[2]int{153, 148}, true},
		{[2]int{150, 150}, false},
		{[2]int{201, 159}, true},
		{[2]int{200, 139}, true},
		{[2]int{200, 199}, true},
		{[2]int{200, 59}, true},
		{[2]int{250, 249}, true},
		{[2]int{251, 149}, true},
		{[2]int{305, 232}, true},
		{[2]int{405, 232}, true},
		{[2]int{455, 448}, true},
		{[2]int{455, 450}, false},
	}
	for _, c := range cases {
		got := TweetWorthy(c.in[0], c.in[1])
		if got != c.want {
			t.Errorf("TweetWorthy(%d, %d) == %q, want %q", c.in[0], c.in[1], got, c.want)
		}
	}
}
