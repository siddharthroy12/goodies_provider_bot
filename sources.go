package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (a *Application) sendPornForSubRedditPag(page string, chatId int64) error {
	html, err := a.fetchHTML(page)

	if err != nil {
		return err
	}

	links, err := extractPreviewReddItLinks(strings.NewReader(html))

	if err != nil {
		return err
	}

	// Send links as media to chat
	if len(links) == 0 {
		return errors.New("no links found in page")
	}
	links = removeDuplicateLinks(links)

	// Send each link as actual downloaded media
	for _, link := range links {
		fmt.Printf("Sending link %s\n", link)

		// Download and send the media file
		err = a.downloadAndSendMedia(chatId, link)
		if err != nil {
			fmt.Printf("Failed to download and send media from %s: %v\n", link, err)
			// Continue with next link even if one fails
		}

		// Small delay to prevent rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (a *Application) sendRedditPorn(chatId int64) error {

	var subreddits = []string{
		"hentai",
		"HENTAI_GIF",
		"jerkbudsHentai",
		"Hentai_AnimeNSFW",
		"HelplessHentai",
		"HentaiM",
		"FreeuseHentai",
		"nhentai",
		"JerkOffToAnime",
		"rule34",
		"PublicHentai",
		"hentainmanga",
		"hentainmanga",
		"Total_Hentai",
		"hentai_irl",
	}
	var err error

	for _, subreddit := range subreddits {
		// First try to send best
		err = a.sendPornForSubRedditPag(fmt.Sprintf("https://www.reddit.com/r/%s", subreddit), chatId)
		if err != nil {
			a.sendText(chatId, fmt.Sprintf("Something went wrong: %s", err.Error()))
		}
		// Then send hot
		err = a.sendPornForSubRedditPag(fmt.Sprintf("https://www.reddit.com/r/%s/hot", subreddit), chatId)
		if err != nil {
			a.sendText(chatId, fmt.Sprintf("Something went wrong: %s", err.Error()))
		}
		// Then send top
		err = a.sendPornForSubRedditPag(fmt.Sprintf("https://www.reddit.com/r/%s/top", subreddit), chatId)
		if err != nil {
			a.sendText(chatId, fmt.Sprintf("Something went wrong: %s", err.Error()))
		}
		// Then send rising
		err = a.sendPornForSubRedditPag(fmt.Sprintf("https://www.reddit.com/r/%s/rising", subreddit), chatId)
		if err != nil {
			a.sendText(chatId, fmt.Sprintf("Something went wrong: %s", err.Error()))
		}
	}

	return err
}
