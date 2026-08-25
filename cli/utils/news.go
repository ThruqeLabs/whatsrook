package cliutils

import (
	"whatsrook/cli/utils/news"
)

// NewsArticle represents a summarized news item.
type NewsArticle = news.NewsArticle

// APNewsArticleDetail represents a fully parsed AP News article.
type APNewsArticleDetail = news.APNewsArticleDetail

// NewsSession tracks the active news query session for a chat.
type NewsSession = news.NewsSession

// WABetaArticle represents a parsed WABetaInfo post.
type WABetaArticle = news.WABetaArticle

var (
	// SetRecentNewsSession caches articles for numbered selection.
	SetRecentNewsSession = news.SetRecentNewsSession

	// GetRecentNewsSession retrieves cached articles for numbered selection.
	GetRecentNewsSession = news.GetRecentNewsSession

	// CleanHTMLText strips tags and unescapes HTML entities.
	CleanHTMLText = news.CleanHTMLText

	// ParseAPNewsHTML parses AP News hub HTML for articles.
	ParseAPNewsHTML = news.ParseAPNewsHTML

	// FetchAPNews fetches news articles for a specific hub/country.
	FetchAPNews = news.FetchAPNews

	// FetchNewsImage downloads image bytes for a news thumbnail.
	FetchNewsImage = news.FetchNewsImage

	// ParseAPNewsArticle extracts title, lead image URL, and formatted body paragraphs from AP News article HTML.
	ParseAPNewsArticle = news.ParseAPNewsArticle

	// FetchAPNewsArticle fetches and parses the full story for an AP News article URL.
	FetchAPNewsArticle = news.FetchAPNewsArticle

	// ExtractAPNewsArticleURLs extracts all unique AP News article URLs from a message text.
	ExtractAPNewsArticleURLs = news.ExtractAPNewsArticleURLs

	// ParseNewsSelectionNumber parses article numbers from user input.
	ParseNewsSelectionNumber = news.ParseNewsSelectionNumber

	// ParseRecentWABetaLink extracts the latest article URL from wabetainfo.com homepage HTML.
	ParseRecentWABetaLink = news.ParseRecentWABetaLink

	// ParseWABetaArticle extracts title, featured/content image URL, and formatted content from article HTML.
	ParseWABetaArticle = news.ParseWABetaArticle

	// FetchWABetaLatest fetches the wabetainfo.com homepage, extracts the latest article URL, and parses the article.
	FetchWABetaLatest = news.FetchWABetaLatest
)
