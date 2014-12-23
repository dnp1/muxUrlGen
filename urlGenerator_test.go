package urlGen

import "testing"
import "sort"
import "log"
import "io/ioutil"

type Cases struct {
	in   string
	want []string
}

// auxiliary function to compare string slices
func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// testing expected errors (from invalid urls)
func TestExpectedErros(t *testing.T) {

	//Discarding logs from application
	log.SetOutput(ioutil.Discard)

	cases := []Cases{
		{
			in:   "/handler/",
			want: []string{"/handler/"}, //irrelevant here
		},
		{
			in: "/handler/{id}/{iqWd: [0-9]}",
		},
		{
			in: "/handler/{id}/{id:[0-9]+}",
		},
	}

	//Test Short mode
	for _, c := range cases {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected error has not ocurred. Case: %q", c)
				}
			}()
			_ = GetUrlVarsPermutations(c.in, false)
		}()
	}
}

// testing urls with no vars, they must generate slices with themselves
func TestNonVarUrl(t *testing.T) {
	cases := []Cases{
		{
			in:   "/handler",
			want: []string{"/handler"},
		},
		{
			in:   "/handler/url/home",
			want: []string{"/handler/url/home"},
		},
		{
			in:   "/handl/u_12%20rl/ho12fme",
			want: []string{"/handl/u_12%20rl/ho12fme"},
		},
	}

	//Test Short mode
	for _, c := range cases {
		got := GetUrlVarsPermutations(c.in, false)
		sort.Strings(got)
		sort.Strings(c.want)

		if !testEq(c.want, got) {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestVarUrlShortMode(t *testing.T) {
	cases := []Cases{
		{ //Simple pattern
			in:   "/handler/{id}",
			want: []string{"/handler/id/{id}"},
		},
		{ //Testing simple permutation
			in: "/handler/{id}/{k_id:[0-9]+}",
			want: []string{
				"/handler/id/{id}/k_id/{k_id:[0-9]+}",
				"/handler/k_id/{k_id:[0-9]+}/id/{id}",
			},
		},
		{ //Testing simple permutation with optional var
			in: "/handler/{id}?/{k_id:[0-9]+}",
			want: []string{
				"/handler/k_id/{k_id:[0-9]+}",
				"/handler/id/{id}/k_id/{k_id:[0-9]+}",
				"/handler/k_id/{k_id:[0-9]+}/id/{id}",
			},
		},
		{ //A bit more complex case
			in: "/handler/{k_id:[0-9]+}/{id}/{a:[0-9A-Za-z]}",
			want: []string{
				"/handler/id/{id}/k_id/{k_id:[0-9]+}/a/{a:[0-9A-Za-z]}",
				"/handler/id/{id}/a/{a:[0-9A-Za-z]}/k_id/{k_id:[0-9]+}",
				"/handler/k_id/{k_id:[0-9]+}/id/{id}/a/{a:[0-9A-Za-z]}",
				"/handler/k_id/{k_id:[0-9]+}/a/{a:[0-9A-Za-z]}/id/{id}",
				"/handler/a/{a:[0-9A-Za-z]}/k_id/{k_id:[0-9]+}/id/{id}",
				"/handler/a/{a:[0-9A-Za-z]}/id/{id}/k_id/{k_id:[0-9]+}",
			},
		},
		{ //A bit more complex case with optional vars
			in: "/handler/{k_id:[0-9]+}/{id}/{a:[0-9A-Za-z]}?",
			want: []string{
				"/handler/id/{id}/k_id/{k_id:[0-9]+}",
				"/handler/k_id/{k_id:[0-9]+}/id/{id}",
				"/handler/id/{id}/k_id/{k_id:[0-9]+}/a/{a:[0-9A-Za-z]}",
				"/handler/id/{id}/a/{a:[0-9A-Za-z]}/k_id/{k_id:[0-9]+}",
				"/handler/k_id/{k_id:[0-9]+}/id/{id}/a/{a:[0-9A-Za-z]}",
				"/handler/k_id/{k_id:[0-9]+}/a/{a:[0-9A-Za-z]}/id/{id}",
				"/handler/a/{a:[0-9A-Za-z]}/k_id/{k_id:[0-9]+}/id/{id}",
				"/handler/a/{a:[0-9A-Za-z]}/id/{id}/k_id/{k_id:[0-9]+}",
			},
		},
	}

	for _, c := range cases {
		got := GetUrlVarsPermutations(c.in, false)
		sort.Strings(got)
		sort.Strings(c.want)

		if !testEq(c.want, got) {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestVarUrlLongMode(t *testing.T) {
	cases := []Cases{
		{ //Simple pattern
			in:   "/handler/id-item/{id}",
			want: []string{"/handler/id-item/{id}"},
		},
		{ //Testing simple permutation
			in: "/handler/id/{id}/id_K/{k_id:[0-9]+}",
			want: []string{
				"/handler/id/{id}/id_K/{k_id:[0-9]+}",
				"/handler/id_K/{k_id:[0-9]+}/id/{id}",
			},
		},
		{ //Testing simple permutation with optional var
			in: "/handler/id/{id}?/id_K/{k_id:[0-9]+}",
			want: []string{
				"/handler/id_K/{k_id:[0-9]+}",
				"/handler/id/{id}/id_K/{k_id:[0-9]+}",
				"/handler/id_K/{k_id:[0-9]+}/id/{id}",
			},
		},
		{ //A bit more complex case
			in: "/handler/id_K/{k_id:[0-9]+}/id/{id}/VAR-a/{a:[0-9A-Za-z]}",
			want: []string{
				"/handler/id/{id}/id_K/{k_id:[0-9]+}/VAR-a/{a:[0-9A-Za-z]}",
				"/handler/id/{id}/VAR-a/{a:[0-9A-Za-z]}/id_K/{k_id:[0-9]+}",
				"/handler/id_K/{k_id:[0-9]+}/id/{id}/VAR-a/{a:[0-9A-Za-z]}",
				"/handler/id_K/{k_id:[0-9]+}/VAR-a/{a:[0-9A-Za-z]}/id/{id}",
				"/handler/VAR-a/{a:[0-9A-Za-z]}/id_K/{k_id:[0-9]+}/id/{id}",
				"/handler/VAR-a/{a:[0-9A-Za-z]}/id/{id}/id_K/{k_id:[0-9]+}",
			},
		},
		{ //A bit more complex case with optional vars
			in: "/handler/id_K/{k_id:[0-9]+}/id/{id}/VAR-a/{a:[0-9A-Za-z]}?",
			want: []string{
				"/handler/id/{id}/id_K/{k_id:[0-9]+}",
				"/handler/id_K/{k_id:[0-9]+}/id/{id}",
				"/handler/id/{id}/id_K/{k_id:[0-9]+}/VAR-a/{a:[0-9A-Za-z]}",
				"/handler/id/{id}/VAR-a/{a:[0-9A-Za-z]}/id_K/{k_id:[0-9]+}",
				"/handler/id_K/{k_id:[0-9]+}/id/{id}/VAR-a/{a:[0-9A-Za-z]}",
				"/handler/id_K/{k_id:[0-9]+}/VAR-a/{a:[0-9A-Za-z]}/id/{id}",
				"/handler/VAR-a/{a:[0-9A-Za-z]}/id_K/{k_id:[0-9]+}/id/{id}",
				"/handler/VAR-a/{a:[0-9A-Za-z]}/id/{id}/id_K/{k_id:[0-9]+}",
			},
		},
	}

	for _, c := range cases {
		got := GetUrlVarsPermutations(c.in, true)
		sort.Strings(got)
		sort.Strings(c.want)

		if !testEq(c.want, got) {
			t.Errorf("Permutations of (%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
