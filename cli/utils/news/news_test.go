package news

import (
	"strings"
	"testing"
)

const sampleArticleHTML = `<!DOCTYPE html>
<html>
<head>
    <meta property="og:title" content="Nigeria orders investigation into second fake government agency | AP News" />
    <meta property="og:image" content="https://dims.apnews.com/dims4/default/12345/image.jpg" />
</head>
<body>
    <h1 class="Page-headline">Nigeria orders investigation into second fake government agency</h1>
    <div class="RichTextStoryBody RichTextBody">
        <p>ABUJA, Nigeria (AP) — <span class="LinkEnhancement"><a class="Link AnClick-LinkEnhancement" data-gtm-enhancement-style="LinkEnhancementA" href="https://apnews.com/hub/nigeria">Nigeria’s</a></span> President Bola Tinubu on Friday ordered an investigation into a man accused of creating and running a fake government agency, the second such agency uncovered in the country in recent weeks.</p>
        <p>Musa Adamu Aliyu, the head of the country’s corruption watchdog, announced the investigation after briefing the president in the capital, Abuja.</p>
        <p>The latest bogus agency, called the National Brands Development and Made in Nigeria Special Project Office, was operating illegally from a government office and was headed by George Buchi Nwabueze, authorities said.</p>
        <p>The scandal is the latest involving an allegedly fake government agency under Tinubu’s administration. Last month, authorities discovered that another purported agency, the Presidential Foreign Investment Promotion Council, had operated for years, with its director meeting government officials and accessing public funds.</p>
        <div id="optimizelyHubpeekId" class="optimizelyHubpeekClass"></div>
        <div class="FreeStar Advertisement" data-module="">
            <div data-freestar-ad="__336x280 __336x280" id="" class="fs-feed-ad">
                <script data-cfasync="false" type="text/javascript">
                    freestar.queue.push(function () {});
                </script>
            </div>
        </div>
        <p>The man accused of creating the first bogus agency was arrested by police. Aliyu said the investigation into the first agency led to the discovery of the second.</p>
        <p>Analysts believe the perpetrators are working with senior government officials to use the agencies to siphon money from the national budget.</p>
        <p>Aliyu did not say whether the accused man had profited from the scheme before it was uncovered.</p>
        <p>Three senior government officials have been suspended, and authorities have issued an arrest warrant for Nwabueze.</p>
        <p>____</p>
        <p>AP’s Africa coverage at: <span class="LinkEnhancement"><a class="Link AnClick-LinkEnhancement" data-gtm-enhancement-style="LinkEnhancementA" href="https://apnews.com/hub/africa">https://apnews.com/hub/africa</a></span></p>
    </div>
</body>
</html>`

func TestParseAPNewsArticle(t *testing.T) {
	articleURL := "https://apnews.com/article/nigeria-fake-government-agency-bola-tinubu-5914f1c917b8dc309b3c6b71e78a1ea1"
	art, err := ParseAPNewsArticle(sampleArticleHTML, articleURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTitle := "Nigeria orders investigation into second fake government agency"
	if art.Title != expectedTitle {
		t.Errorf("expected Title %q, got %q", expectedTitle, art.Title)
	}

	if art.URL != articleURL {
		t.Errorf("expected URL %q, got %q", articleURL, art.URL)
	}

	if art.ImageURL != "https://dims.apnews.com/dims4/default/12345/image.jpg" {
		t.Errorf("expected ImageURL to be parsed, got %q", art.ImageURL)
	}

	if !strings.Contains(art.Content, "ABUJA, Nigeria (AP) — Nigeria’s President Bola Tinubu") {
		t.Errorf("expected content to contain first paragraph, got %q", art.Content)
	}

	if strings.Contains(art.Content, "AP’s Africa coverage at:") {
		t.Errorf("expected boilerplate coverage text to be stripped, but got: %q", art.Content)
	}

	if strings.Contains(art.Content, "____") {
		t.Errorf("expected divider underscores to be stripped, but got: %q", art.Content)
	}
}

func TestExtractAPNewsArticleURLs(t *testing.T) {
	input := "Check this story: https://apnews.com/article/world-news-123 and also https://apnews.com/article/us-politics-456."
	urls := ExtractAPNewsArticleURLs(input)
	if len(urls) != 2 {
		t.Fatalf("expected 2 URLs, got %d", len(urls))
	}
	if urls[0] != "https://apnews.com/article/world-news-123" {
		t.Errorf("expected first url %q, got %q", "https://apnews.com/article/world-news-123", urls[0])
	}
	if urls[1] != "https://apnews.com/article/us-politics-456" {
		t.Errorf("expected second url %q, got %q", "https://apnews.com/article/us-politics-456", urls[1])
	}
}

func TestParseNewsSelectionNumber(t *testing.T) {
	cases := []struct {
		input string
		want  int
		ok    bool
	}{
		{"1", 1, true},
		{"#2", 2, true},
		{"article 3", 3, true},
		{".news 4", 4, true},
		{"news 5", 5, true},
		{"invalid", 0, false},
		{"0", 0, false},
		{"55", 0, false},
	}

	for _, c := range cases {
		got, ok := ParseNewsSelectionNumber(c.input)
		if ok != c.ok || got != c.want {
			t.Errorf("ParseNewsSelectionNumber(%q) = (%d, %v); want (%d, %v)", c.input, got, ok, c.want, c.ok)
		}
	}
}
