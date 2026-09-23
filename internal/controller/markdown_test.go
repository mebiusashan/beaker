package controller

import (
	"strings"
	"testing"
)

func TestRenderMarkdownPreservesCommonExtensionsAndFiltersRawHTML(t *testing.T) {
	markdown := "```go\nfmt.Println(\"hello\")\n```\n\n" +
		"| Name | Value |\n| --- | --- |\n| mode | fast |\n\n" +
		"<script>alert(1)</script>"

	html := renderMarkdown(markdown)

	for _, want := range []string{
		"<pre><code",
		`fmt.Println(&quot;hello&quot;)`,
		"<table>",
		"<th>Name</th>",
		"<td>fast</td>",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered Markdown missing %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, "<script") {
		t.Errorf("raw script HTML was not filtered:\n%s", html)
	}
}
