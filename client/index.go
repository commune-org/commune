package client

import (
	"commune/gomatrix"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/unrolled/secure"
)

type Room struct {
	Path  string `json:"path"`
	Alias string `json:"alias"`
}

func (c *Client) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		/*
			c.IndexUser(w, r)
			return
		*/

		us := LoggedInUser(r)

		if us != nil {
			c.IndexUser(w, r)
			return
		}

		type page struct {
			BasePage
			LoginError bool
		}

		t := page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("login-error")
			if len(x) > 0 {
				t.LoginError = true
			}
			s.Save(r, w)
		}

		nonce := secure.CSPNonce(r.Context())

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "index", t)
	}
}

func (c *Client) IndexUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Accept-Ranges", "bytes")

	user := LoggedInUser(r)

	type page struct {
		Room Room
		BasePage
		LoggedInUser     *User
		State            interface{}
		ProfileLink      template.URL
		Posts            interface{}
		HomeServerURL    string
		FederationDomain string
		ShortlinkDomain  string
		FeedItems        interface{}
		InitialPosts     interface{}
		Depth            int
		TenorKey         string
		NewUser          bool
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{
		Room: Room{
			Path: "user-index",
		},
	}

	s, err := GetSession(r, c)
	if err != nil {
		log.Println(err)
	}
	if s != nil {
		x := s.Flashes("signed-up")
		if len(x) > 0 {
			t.NewUser = true
		}
		s.Save(r, w)
	}

	if user != nil {

		acc, err := c.ConstructAccount(user)

		if err != nil {
			c.Error(w, r)
			return
		}

		us := acc.User

		for i, altUser := range us.AltAccounts {
			ac, err := c.ConstructAccount(altUser)

			if err != nil {
				c.Error(w, r)
				return
			}

			au := ac.User

			au.SyncState = ac.SyncState
			au.Spaces = ac.Spaces
			au.TimelineState = ac.TimelineState
			au.Notifications = ac.Notifications
			au.UserStatus = ac.UserStatus
			au.AccountData = ac.AccountData

			us.AltAccounts[i] = &au
		}

		us.SyncState = acc.SyncState
		us.Spaces = acc.Spaces
		us.TimelineState = acc.TimelineState
		us.Notifications = acc.Notifications
		us.UserStatus = acc.UserStatus
		us.AccountData = acc.AccountData
		us.Profile = acc.Profile

		t.LoggedInUser = &us
	}

	if c.Config.Mode == "development" {
		t.HomeServerURL = c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)
		t.ShortlinkDomain = c.Config.Client.ShortlinkDomain + fmt.Sprintf(`:%d`, c.Config.Client.Port)
	} else {
		t.HomeServerURL = fmt.Sprintf(`https://%s`, c.Config.Matrix.PublicServer)
		t.ShortlinkDomain = c.Config.Client.ShortlinkDomain
	}
	t.FederationDomain = c.Config.Matrix.FederationServer

	t.Nonce = nonce

	t.TenorKey = c.Config.Tenor.Key

	c.Templates.ExecuteTemplate(w, "index-user", t)
}

type ConstructAccountResponse struct {
	SyncState     interface{}
	Notifications interface{}
	AccountData   interface{}
	UserStatus    interface{}
	Spaces        interface{}
	TimelineState interface{}
	Profile       interface{}
	User          User
}

func (c *Client) ConstructAccount(user *User) (ConstructAccountResponse, error) {

	resp := ConstructAccountResponse{}

	serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

	matrix, err := gomatrix.NewClient(serverName, user.UserID, user.MatrixAccessToken)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	rms, err := c.SyncUserState(matrix)
	if err != nil {
		log.Println(err)
	}

	if rms != nil {
		resp.SyncState = rms

		type Space struct {
			RoomID         string   `json:"room_id,omitempty"`
			Name           string   `json:"name,omitempty"`
			Alias          string   `json:"alias,omitempty"`
			Pathname       string   `json:"pathname,omitempty"`
			CreatedOn      int64    `json:"created_on,omitempty"`
			OriginServerTS int64    `json:"origin_server_ts,omitempty"`
			Sender         string   `json:"sender,omitempty"`
			Topic          string   `json:"topic,omitempty"`
			Owner          bool     `json:"owner,omitempty"`
			Avatar         string   `json:"avatar,omitempty"`
			RoomType       string   `json:"room_type,omitempty"`
			Child          bool     `json:"child,omitempty"`
			SpaceID        string   `json:"space_id,omitempty"`
			SpacePathname  string   `json:"space_pathname,omitempty"`
			Rooms          []*Space `json:"rooms,omitempty"`
		}

		constructed := []Space{}
		spaces := []Space{}

		var buildSpace func(events []gomatrix.Event, spaceID string) Space

		buildSpace = func(events []gomatrix.Event, spaceID string) Space {
			space := Space{}
			for _, event := range events {
				if event.Type == "m.space.parent" {
					space.Child = true
					space.SpaceID = *event.StateKey
				}
				if event.Type == "m.room.create" {
					space.Sender = event.Sender
					space.OriginServerTS = event.Timestamp
					if event.Sender == user.UserID {
						space.Owner = true
					}
				}
				if event.Type == "commune.room" {
					if rt, ok := event.Content["room_type"].(string); ok {
						space.RoomType = rt
					}
				}
				if event.Type == "m.room.name" {
					if name, ok := event.Content["name"].(string); ok {
						space.Name = name
					}
				}
				if event.Type == "m.room.canonical_alias" {
					if alias, ok := event.Content["alias"].(string); ok {
						x := strings.Split(alias, ":")[0]
						a := x[1:]
						space.Alias = a
						space.Pathname = fmt.Sprintf(`/%s`, a)
					}
				}
				if event.Type == "m.room.topic" {
					if topic, ok := event.Content["topic"].(string); ok {
						space.Topic = topic
					}
				}
				if event.Type == "m.room.avatar" {
					if avatar, ok := event.Content["url"].(string); ok {
						space.Avatar = avatar
					}
				}
				if event.Type == "m.room.topic" {
					if topic, ok := event.Content["topic"].(string); ok {
						space.Topic = topic
					}
				}
				if event.Type == "commune.room" {
					if rt, ok := event.Content["room_Type"].(string); ok {
						space.RoomType = rt
					}
				}
				if event.Type == "m.space.child" {
					roomID := *event.StateKey
					for id, rm := range rms.(*gomatrix.RespSync).Rooms.Join {
						if id == roomID {
							child := buildSpace(rm.State.Events, "")
							child.RoomID = id
							child.SpaceID = spaceID
							space.Rooms = append(space.Rooms, &child)
						}
					}
				}
			}

			for _, child := range space.Rooms {
				child.SpacePathname = space.Pathname
			}

			return space

		}

		buildSpaces := func() []Space {

			for roomID, room := range rms.(*gomatrix.RespSync).Rooms.Join {
				space := buildSpace(room.State.Events, roomID)
				space.RoomID = roomID
				spaces = append(spaces, space)
			}

			return spaces
		}

		tempSpaces := buildSpaces()

		for _, space := range tempSpaces {
			if !space.Child {
				constructed = append(constructed, space)
			}
		}

		resp.Spaces = constructed

	}

	tms, err := c.SyncUserTimeline(matrix)
	if err != nil {
		log.Println(err)
	}
	if tms != nil {
		resp.TimelineState = tms
	}

	notifications, err := matrix.GetNotifications("")
	if err != nil {
		log.Println(err)
	}
	if notifications != nil {
		resp.Notifications = notifications
	}

	status, err := matrix.GetStatus(user.UserID)
	if err != nil {
		log.Println(err)
	}

	if status != nil {
		resp.UserStatus = status
	}

	profile, err := matrix.GetProfile(user.UserID)
	if err != nil {
		log.Println(err)
	}

	if profile != nil {
		if profile.Displayname != nil && len(*profile.Displayname) > 0 {
			user.DisplayName = *profile.Displayname
		}
		if profile.AvatarURL != nil && len(*profile.AvatarURL) > 0 {
			/*
				av := c.BuildDownloadLink(*profile.AvatarURL)
				user.AvatarURL = StripMXCPrefix(av)
			*/
			user.AvatarURL = *profile.AvatarURL
		}
	}

	username := GetLocalPart(user.UserID)

	ds := fmt.Sprintf(`#@%s:%s`, username, c.Config.Matrix.FederationServer)

	ra, err := matrix.ResolveAlias(ds)
	if err != nil {
		log.Println(err)
		log.Println("profile room missing?")
	}

	prs, err := matrix.RoomState(ra.RoomID.String())
	if err != nil {
		log.Println(err)
		log.Println("profile room missing?")
	}

	if err != nil {
		log.Println(err)
	}

	type Profile struct {
		Banner string `json:"banner"`
		About  string `json:"about"`
		Color  string `json:"color"`
	}

	pr := Profile{}

	if prs != nil && len(prs) > 0 {
		for _, event := range prs {

			if event.Type == "commune.profile.banner" {
				if b, ok := event.Content["url"].(string); ok {
					pr.Banner = b
				}
			}
			if event.Type == "commune.profile.about" {
				if b, ok := event.Content["about"].(string); ok {
					pr.About = b
				}
			}
			if event.Type == "commune.profile.color" {
				if b, ok := event.Content["color"].(string); ok {
					pr.Color = b
				}
			}
		}
	}

	resp.Profile = pr

	data, err := matrix.GetAccountData(user.UserID, "m.direct")
	if err != nil {
		log.Println(err)
	}

	if data != nil {
		resp.AccountData = data
	}

	resp.User = *user

	return resp, nil
}
