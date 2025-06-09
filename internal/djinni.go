package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/vmamchur/joblin-scraper/db/generated"
)

const baseUrl = "https://djinni.co"

type DjinniScraper struct {
	q        *generated.Queries
	email    string
	password string
}

func (d DjinniScraper) Scrape() error {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), "http://chrome:9222/json/version")
	defer allocCancel()

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	var currentUrl string
	err := chromedp.Run(ctx,
		chromedp.Navigate(baseUrl+"/login"),
		chromedp.Location(&currentUrl),
	)
	if err != nil {
		return err
	}

	if strings.Contains(currentUrl, "/login") {
		err = chromedp.Run(ctx,
			chromedp.WaitVisible("form#signup", chromedp.ByQuery),
			chromedp.SendKeys(`form#signup input[name="email"]`, d.email, chromedp.ByQuery),
			chromedp.SendKeys(`form#signup input[name="password"]`, d.password, chromedp.ByQuery),
			chromedp.Click(`form#signup button[type="submit"]`, chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			return err
		}
	}

	page := 1

	for {
		log.Printf("Navigating to page %d...", page)

		var jobNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf(baseUrl+"/jobs/?primary_keyword=fullstack&page=%d", page)),
			chromedp.WaitVisible("li[id^=job-item-]", chromedp.ByQuery),
			chromedp.Nodes("li[id^=job-item-]", &jobNodes, chromedp.ByQueryAll),
		)
		if err != nil {
			return err
		}

		log.Printf("Found %d vacancies on page %d", len(jobNodes), page)

		for _, node := range jobNodes {
			var title, url, companyName string

			err := chromedp.Run(ctx,
				chromedp.Text(".job-item__title-link", &title, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.AttributeValue(".job-item__title-link", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.Text(`a[data-analytics="company_page"]`, &companyName, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(node)),
			)
			if err != nil {
				log.Printf("Skipping vacancy due to extraction error: %v", err)
				continue
			}

			fullUrl := baseUrl + url
			log.Printf("Processing: %s", fullUrl)

			_, err = d.q.CreateVacancy(ctx, generated.CreateVacancyParams{
				Title:       title,
				CompanyName: companyName,
				Url:         fullUrl,
			})
			if err != nil {
				log.Printf("Error saving %s: %v", fullUrl, err)
				continue
			}

			log.Printf("Saved: %s", fullUrl)
		}

		var isNextBtnVisible bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('li.page-item:not(.disabled) a.page-link span.bi-chevron-right') !== null`, &isNextBtnVisible),
		)
		if err != nil {
			return err
		}
		if !isNextBtnVisible {
			log.Println("No more pages. Done scraping.")
			break
		}

		page++
	}

	return nil
}
