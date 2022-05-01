package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (c *Client) LinkMetadata() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Href string `json:"href"`
		}

		var pay payload
		if r.Body == nil {
			log.Println(err)
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Title            string  `json:"title,omitempty"`
			Author           string  `json:"author,omitempty"`
			Description      string  `json:"description,omitempty"`
			Image            string  `json:"image,omitempty"`
			IsYoutube        *bool   `json:"is_youtube,omitempty"`
			YoutubeID        *string `json:"youtube_id,omitempty"`
			IsVimeo          *bool   `json:"is_vimeo,omitempty"`
			VimeoID          *string `json:"vimeo_id,omitempty"`
			SoundCloudPlayer *string `json:"sound_cloud_player,omitempty"`
			IsWikipedia      *bool   `json:"is_wikipedia,omitempty"`
		}

		ff := Response{}

		up, err := url.Parse(pay.Href)
		if err != nil || up == nil || up.Host == "" {
			http.Error(w, err.Error(), 400)
			return
		}

		key := c.Config.YoutubeKey

		isYoutube := up.Host == "www.youtube.com" || up.Host == "youtube.com"
		isShortYoutube := up.Host == "youtu.be"

		isVimeo := up.Host == "www.vimeo.com" || up.Host == "vimeo.com"

		isSoundCloud := up.Host == "soundcloud.com" || up.Host == "www.soundcloud.com" || up.Host == "m.soundcloud.com"

		isWikipedia := up.Host == "en.wikipedia.org"

		var title, description, image, author string

		if (isYoutube || isShortYoutube) && len(key) > 0 {

			m, err := url.ParseQuery(up.RawQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			id := ""

			if isYoutube {

				if len(m) > 0 {
					id = m["v"][0]
				}
			} else if isShortYoutube {

				if isShortYoutube {
					id = up.Path[1:]
				}
			}

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			service, err := youtube.NewService(ctx, option.WithAPIKey(key))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			videos := service.Videos.List([]string{"id", "snippet"})

			videos = videos.Id(id)

			response, err := videos.Do()

			if err != nil {
				log.Println(err)
			}

			if response != nil {

				items := response.Items

				if items != nil && len(items) >= 1 {
					title = items[0].Snippet.Title
					description = items[0].Snippet.Description
				}

				if len(description) > 500 {
					description = description[:500]
				}

				yt := true
				ff.IsYoutube = &yt
				ff.YoutubeID = &id

			}

		} else if isWikipedia {
			path := strings.TrimLeft(up.Path, "/")
			spl := strings.Split(path, "/")
			log.Println(spl)

			w, err := mwclient.New("https://en.wikipedia.org/w/api.php", "communeWikiScraper")
			if err != nil {
				log.Println(err)
			}

			parameters := params.Values{
				"action":    "query",
				"prop":      "extracts",
				"format":    "html",
				"exintro":   "",
				"redirects": "1",
				"titles":    spl[1],
			}
			response, err := w.Get(parameters)
			if err != nil {
				log.Println(err)
			}

			pages, err := response.GetObjectArray("query", "pages")
			if err != nil {
				log.Println(err)
			}

			var extract, ti string
			for _, item := range pages {
				extract, err = item.GetString("extract")
				if err != nil {
					log.Println(err)
				}
				ti, err = item.GetString("title")
				if err != nil {
					log.Println(err)
				}
			}

			title = ti
			description = extract
			wi := true
			ff.IsWikipedia = &wi
		} else if isVimeo {
			type video struct {
				Type         string `json:"type"`
				Title        string `json:"title"`
				HTML         string `json:"html"`
				Description  string `json:"description"`
				ThumbnailURL string `json:"thumbnail_url"`
				VideoID      int64  `json:"video_id"`
			}

			id := up.Path[1:]

			endpoint := fmt.Sprintf(`https://vimeo.com/api/oembed.json?url=https://vimeo.com/%s`, id)

			resp, err := http.Get(endpoint)
			if err != nil {
				log.Fatalln(err)
			}

			defer resp.Body.Close()
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			var vid video
			json.Unmarshal(bodyBytes, &vid)

			desc := vid.Description
			if len(desc) > 160 {
				desc = desc[:157]
			}

			title = vid.Title
			description = desc
			image = vid.ThumbnailURL

			if vid.Type == "video" {
				vimeoID := strconv.FormatInt(vid.VideoID, 10)
				vi := true
				ff.IsVimeo = &vi
				ff.VimeoID = &vimeoID
			}

		} else {
			md := c.Scrape(pay.Href, up.Host)

			title = md.Title
			description = md.Description
			author = md.Author
			image = md.Image

			if isSoundCloud && md.SoundCloud != "" {
				scPlayer := md.SoundCloud
				scPlayer = strings.ReplaceAll(scPlayer, `auto_play=false`, `auto_play=true`)

				ff.SoundCloudPlayer = &scPlayer

				image = strings.ReplaceAll(image, "500x500", "200x200")
			}
		}

		ff.Title = title
		ff.Author = author
		ff.Description = description
		ff.Image = image

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

type LinkMetaData struct {
	Title       string
	Description string
	Image       string
	Author      string
	SoundCloud  string
}

func (c *Client) Scrape(link string, domain string) LinkMetaData {
	co := colly.NewCollector()
	extensions.RandomUserAgent(co)
	extensions.Referer(co)

	lmd := LinkMetaData{}

	isSoundCloud := domain == "soundcloud.com" || domain == "www.soundcloud.com"

	co.OnHTML("head", func(e *colly.HTMLElement) {

		// Extract meta tags from the document
		metaTags := e.DOM.ParentsUntil("~").Find("meta")

		metaTags.Each(func(_ int, s *goquery.Selection) {
			// Search for og:type meta tags
			name, _ := s.Attr("name")
			prop, _ := s.Attr("property")

			if strings.EqualFold(name, "description") {
				description, _ := s.Attr("content")
				lmd.Description = description
			}

			if strings.EqualFold(name, "author") {
				author, _ := s.Attr("content")
				lmd.Author = author
			}

			if strings.EqualFold(prop, "og:image") {
				image, _ := s.Attr("content")
				lmd.Image = image
			}

			if lmd.Image == "" && strings.EqualFold(prop, "twitter:image:src") {
				image, _ := s.Attr("content")
				lmd.Image = image
			}

			if isSoundCloud && strings.EqualFold(prop, "twitter:player") {
				player, _ := s.Attr("content")
				lmd.SoundCloud = player
			}

		})

	})

	co.OnHTML("head title", func(e *colly.HTMLElement) {
		lmd.Title = e.Text

	})

	co.OnRequest(func(r *colly.Request) {
	})

	co.Visit(link)

	return lmd
}
