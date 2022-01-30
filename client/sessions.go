package client

import (
	"commune/gomatrix"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/tidwall/gjson"
)

type JoinedRoom struct {
	RoomID    string `json:"room_id"`
	RoomAlias string `json:"room_alias"`
	Avatar    string `json:"avatar"`
	Name      string `json:"name"`
}

type OwnedRoom struct {
	RoomID    string `json:"room_id"`
	RoomAlias string `json:"room_alias"`
}

type User struct {
	DisplayName       string                   `json:"display_name"`
	AvatarURL         string                   `json:"avatar_url"`
	AccessToken       string                   `json:"access_token"`
	MatrixAccessToken string                   `json:"matrix_access_token"`
	DeviceID          string                   `json:"device_id"`
	HomeServer        string                   `json:"home_server"`
	UserID            string                   `json:"user_id"`
	RefreshToken      string                   `json:"refresh_token"`
	RoomID            string                   `json:"room_id"`
	JoinedRooms       []JoinedRoom             `json:"joined_rooms"`
	OwnedRooms        []JoinedRoom             `json:"owned_rooms"`
	WellKnown         string                   `json:"well_known"`
	Federated         bool                     `json:"federated"`
	Preferences       gomatrix.UserPreferences `json:"preferences"`
	Email             string                   `json:"email"`
	EmailVerified     bool                     `json:"email_verified"`
	AltAccounts       []*User                  `json:"alt_accounts"`
	SyncState         interface{}              `json:"sync_state"`
	AccountData       interface{}              `json:"account_data"`
	Notifications     interface{}              `json:"notifications"`
	UserStatus        interface{}              `json:"user_status"`
	TimelineState     interface{}              `json:"timeline_state"`
	Spaces            interface{}              `json:"spaces"`
	Profile           interface{}              `json:"profile"`
}

func NewSession(sec string) *sessions.CookieStore {
	s := sessions.NewCookieStore([]byte(sec))
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: false,
	}
	return s
}

func GetSession(r *http.Request, c *Client) (*sessions.Session, error) {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return s, nil
}

func (c *Client) GetTokenUser(token string) (*User, error) {

	userid, err := c.Store.Get(token).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err := c.Store.Get(userid).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var us User
	err = json.Unmarshal([]byte(user), &us)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &us, nil

}

func (c *Client) AddJoinedRoom(j JoinedRoom, r *http.Request) error {

	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.JoinedRooms = append(us.JoinedRooms, j)

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) UpdateJoinedRooms(matrix *gomatrix.Client, r *http.Request) error {
	rms, err := c.GetUserJoinedRooms(matrix)
	if err != nil {
		log.Println(err)
		return err
	}

	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.JoinedRooms = rms

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) RefreshJoinedRooms(r *http.Request, rooms []JoinedRoom) error {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.JoinedRooms = rooms

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) UpdateUserRoomID(r *http.Request, roomID string) error {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.RoomID = roomID

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) IsRoomMember(r *http.Request, roomID string) (bool, error) {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return false, err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return false, err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return false, err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return false, err
		}

		for i, _ := range us.JoinedRooms {
			for x, _ := range us.JoinedRooms {
				if us.JoinedRooms[i] == us.JoinedRooms[x] {
					return true, nil
				}
			}
		}

	}

	return false, nil
}

func (c *Client) UpdateDisplayName(r *http.Request, displayName string) error {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.DisplayName = displayName

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) SetUserAvatar(r *http.Request, avatar string) error {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		us.AvatarURL = avatar

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (c *Client) SyncUserState(matrix *gomatrix.Client) (interface{}, error) {

	fil := `{
		"room": {
			"timeline": {
				"limit": 0
			}
		}
	}`
	/*
		fil := `{
			"room": {
				"timeline": {
					"limit": 50,
					"types": ["m.room.message", "m.reaction"]
				}
			}
		}`
	*/

	sre, err := matrix.SyncRequest(0, "", fil, true, "offline")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return sre, nil
}

func (c *Client) SyncUserTimeline(matrix *gomatrix.Client) (interface{}, error) {

	fil := `{
		"room": {
			"state": {
				"limit": 0
			},
			"timeline": {
				"limit": 50,
				"types": ["m.room.message", "m.reaction"]
			}
		}
	}`

	sre, err := matrix.SyncRequest(0, "", fil, false, "offline")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return sre, nil
}

func (c *Client) GetUserJoinedRooms(matrix *gomatrix.Client) ([]JoinedRoom, error) {

	fil, err := matrix.CreateFilter([]byte(`
{
	"room": {
		"timeline": {
			"limit": 100,
		}
	}
}
	`))

	sre, err := matrix.SyncRequest(0, "", fil.FilterID, true, "offline")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rms := []JoinedRoom{}
	for roomID, room := range sre.Rooms.Join {

		st, _ := json.Marshal(room.State.Events)
		roomType := gjson.Parse(string(st)).Get(`#(type="commune.room")`).Get("content.room_type")

		if roomType.String() == "page" || roomType.String() == "post" {
			continue
		}
		alias := gjson.Parse(string(st)).Get(`#(type="m.room.canonical_alias")`).Get("content.alias")
		if len(alias.String()) > 0 &&
			!strings.Contains(alias.String(), "#thread") &&
			!strings.Contains(alias.String(), "#public") {

			rm := JoinedRoom{
				RoomID:    roomID,
				RoomAlias: alias.String(),
			}

			avatar := gjson.Parse(string(st)).Get(`#(type="m.room.avatar")`).Get("content.url")
			if avatar.String() != "" {
				rm.Avatar = c.BuildAvatar(avatar.String())
			}

			name := gjson.Parse(string(st)).Get(`#(type="m.room.name")`).Get("content.name")
			if name.String() != "" {
				rm.Name = name.String()
			}

			rms = append(rms, rm)
		}
	}
	sort.Slice(rms, func(i, j int) bool { return rms[i].RoomAlias < rms[j].RoomAlias })
	return rms, nil
}

func (c *Client) RefreshPreferences(r *http.Request) error {
	s, err := c.Sessions.Get(r, c.Config.Client.CookieName)
	if err != nil {
		log.Println(err)
		return err
	}

	token, ok := s.Values["access_token"].(string)
	if ok {
		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			return err
		}
		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			return err
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil {
			log.Println(err)
			return err
		}

		matrix, err := c.TempMatrixClient(us.UserID, us.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			log.Println(err)
			return err
		}

		prefs, err := matrix.GetAccountPreferences(us.UserID)
		if err != nil {
			log.Println(err)
		}

		us.Preferences = *prefs

		serialized, err := json.Marshal(us)
		if err != nil {
			log.Println(err)
			return err
		}

		err = c.Store.Set(userid, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}
