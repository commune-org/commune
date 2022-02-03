package client

import (
	"commune/gomatrix"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (c *Client) Setup() {

	matrix, err := c.TempMatrixClient(c.DefaultUser.UserID, c.DefaultUser.AccessToken)
	if err != nil {
		log.Println(err)
	}

	//fl := genSonyflake()
	//username := strconv.FormatUint(fl, 10)

	//create the room
	pl := gomatrix.Event{
		Type: "m.room.power_levels",
		Content: map[string]interface{}{
			"ban": 60,
			"events": map[string]interface{}{
				"m.room.name":         60,
				"m.room.power_levels": 100,
				"m.room.create":       10,
				"m.space.child":       10,
				"m.space.parent":      10,
			},
			"events_default": 10,
			"invite":         10,
			"kick":           60,
			"notifications": map[string]interface{}{
				"room": 20,
			},
			"redact":        10,
			"state_default": 10,
			"users": map[string]interface{}{
				c.DefaultUser.UserID: 100,
			},
			"users_default": 10,
		},
	}

	initState := []gomatrix.Event{
		gomatrix.Event{
			Type: "m.room.history_visibility",
			Content: map[string]interface{}{
				"history_visibility": "joined",
			},
		}, gomatrix.Event{
			Type: "m.room.guest_access",
			Content: map[string]interface{}{
				"guest_access": "forbidden",
			},
		}, gomatrix.Event{
			Type: "commune.room",
			Content: map[string]interface{}{
				"room_type": "space",
			},
		}, gomatrix.Event{
			Type: "m.room.type",
			Content: map[string]interface{}{
				"type":  "m.space",
				"alias": "commune",
			},
		},
		pl,
	}

	req := &gomatrix.ReqCreateRoom{
		Preset:        "public_chat",
		Visibility:    "public",
		RoomAliasName: "commune",
		Name:          "Commune",
		Topic:         "Commune space.",
		CreationContent: map[string]interface{}{
			"m.federate": true,
		},
		InitialState: initState,
	}

	crr, err := matrix.CreateRoom(req)

	if err != nil || crr == nil {
		log.Println(err)
		log.Println(err)
		return
	}

	//create #general channel

	/*
		f := genSonyflake()
		u := strconv.FormatUint(f, 10)

		initState = []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": "world_readable",
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "commune.room",
				Content: map[string]interface{}{
					"room_type": "board",
				},
			}, gomatrix.Event{
				Type: "m.room.type",
				Content: map[string]interface{}{
					"type":  "m.space",
					"alias": u,
				},
			},
			pl,
		}

		content := map[string]interface{}{
			"canonical_alias": fmt.Sprintf(`#%s:%s`, u, c.Config.Matrix.FederationServer),
			"local_part":      u,
			"via":             []string{c.Config.Client.Domain},
		}

		initState = append(initState, gomatrix.Event{
			Type:     `m.space.parent`,
			StateKey: &crr.RoomID,
			Content:  content,
		})

		req = &gomatrix.ReqCreateRoom{
			Preset:        "public_chat",
			Visibility:    "public",
			Name:          `general`,
			RoomAliasName: u,
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: initState,
		}

		cr, err := matrix.CreateRoom(req)

		if err != nil || cr == nil {
			log.Println(err)
			log.Println(err)
		}

		go func() {
			content["room_id"] = cr.RoomID
			content["room_type"] = "board"
			content["name"] = `general`
			content["default"] = true
			_, err := matrix.SendStateEvent(crr.RoomID, "m.space.child", cr.RoomID, content)
			if err != nil {
				log.Println(err)
			}
		}()
	*/
	type createStreamsRequest struct {
		RoomType string
		ParentID string
	}

	createStream := func(r createStreamsRequest) (string, error) {

		initState := []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": "world_readable",
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "commune.room",
				Content: map[string]interface{}{
					"room_type": r.RoomType,
				},
			},
			pl,
		}

		f := genSonyflake()
		u := strconv.FormatUint(f, 10)

		content := map[string]interface{}{
			"canonical_alias": fmt.Sprintf(`#%s:%s`, u, c.Config.Matrix.FederationServer),
			"local_part":      u,
			"via":             []string{c.Config.Client.Domain},
		}

		initState = append(initState, gomatrix.Event{
			Type:     `m.space.parent`,
			StateKey: &r.ParentID,
			Content:  content,
		})

		req := &gomatrix.ReqCreateRoom{
			Preset:        "public_chat",
			Visibility:    "public",
			RoomAliasName: u,
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: initState,
		}

		cr, err := matrix.CreateRoom(req)

		if err != nil || cr == nil {
			log.Println(err)
			return "", err
		}

		{
			content["room_id"] = cr.RoomID
			content["room_type"] = r.RoomType
			content["default"] = true
			_, err := matrix.SendStateEvent(r.ParentID, "m.space.child", cr.RoomID, content)
			if err != nil {
				log.Println(err)
				return "", err
			}
		}

		return cr.RoomID, nil

	}

	type createChannelRequest struct {
		Name          string
		DefaultStream string
	}

	type streams struct {
		Chat   string
		Topics string
	}

	type createChannelResponse struct {
		RoomID  string
		Streams streams
	}

	createChannel := func(r createChannelRequest) (createChannelResponse, error) {

		//create te general chatroom
		s := genSonyflake()
		v := strconv.FormatUint(s, 10)

		newState := []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": "world_readable",
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "commune.stream",
				Content: map[string]interface{}{
					"default_stream": r.DefaultStream,
				},
			}, gomatrix.Event{
				Type: "m.room.type",
				Content: map[string]interface{}{
					"type":  "m.space",
					"alias": v,
				},
			},
			pl,
		}

		ncontent := map[string]interface{}{
			"canonical_alias": fmt.Sprintf(`#%s:%s`, v, c.Config.Matrix.FederationServer),
			"local_part":      v,
			"via":             []string{c.Config.Client.Domain},
		}
		newState = append(newState, gomatrix.Event{
			Type:     `m.space.parent`,
			StateKey: &crr.RoomID,
			Content:  ncontent,
		})

		requ := &gomatrix.ReqCreateRoom{
			Preset:        "public_chat",
			Visibility:    "public",
			Name:          r.Name,
			RoomAliasName: v,
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: newState,
		}

		ccr, err := matrix.CreateRoom(requ)

		if err != nil || ccr == nil {
			log.Println(err)
			return createChannelResponse{}, nil
		}

		log.Println("CREATE ROOOOOOOOOM")

		rid1, err := createStream(createStreamsRequest{
			RoomType: "chat",
			ParentID: ccr.RoomID,
		})
		if err != nil {
			log.Println(err)
			return createChannelResponse{}, nil
		}
		rid2, err := createStream(createStreamsRequest{
			RoomType: "topics",
			ParentID: ccr.RoomID,
		})
		if err != nil {
			log.Println(err)
			return createChannelResponse{}, nil
		}

		go func() {
			ncontent["room_id"] = ccr.RoomID
			ncontent["name"] = r.Name
			ncontent["default_stream"] = r.DefaultStream
			ncontent["streams"] = map[string]string{
				"chat":   rid1,
				"topics": rid2,
			}
			_, err := matrix.SendStateEvent(crr.RoomID, "m.space.child", ccr.RoomID, ncontent)
			if err != nil {
				log.Println(err)
			}
		}()
		return createChannelResponse{
			RoomID: ccr.RoomID,
			Streams: streams{
				Chat:   rid1,
				Topics: rid2,
			},
		}, nil
	}

	_, err = createChannel(createChannelRequest{
		Name:          "general",
		DefaultStream: "chat",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := createChannel(createChannelRequest{
		Name:          "icebreakers",
		DefaultStream: "topics",
	})
	if err != nil {
		log.Fatal(err)
	}

	{

		topics := []string{
			"Where would you time-travel, if it were possible?",
			"If you could have a superpower, what would it be and why?",
			"If you were the captain of a pirate ship, what would be the name of your ship?",
			"If you could travel to any other planet (real or fictional), where would you go and why?",
			"What is your favorite movie?",
			"Who is your favorite actor/actress?",
			"What is your favorite book?",
			"What is your favorite tv series?",
			"What type of music do you listen to?",
		}

		for _, topic := range topics {

			//figure out why synapse throttles admin room creation
			time.Sleep(1 * time.Second)

			// CREATE TOPIC ROOM
			users := map[string]interface{}{
				c.DefaultUser.UserID: 100,
			}

			pl := gomatrix.Event{
				Type: "m.room.power_levels",
				Content: map[string]interface{}{
					"ban": 60,
					"events": map[string]interface{}{
						"m.room.name":         60,
						"m.room.power_levels": 100,
						"m.room.create":       10,
						"m.space.child":       10,
						"m.space.parent":      10,
					},
					"events_default": 10,
					"invite":         10,
					"kick":           60,
					"notifications": map[string]interface{}{
						"room": 20,
					},
					"redact":        10,
					"state_default": 10,
					"users":         users,
					"users_default": 10,
				},
			}

			initState := []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "commune.room",
					Content: map[string]interface{}{
						"room_type": "space",
					},
				},
				pl,
			}

			req := &gomatrix.ReqCreateRoom{
				Preset:     "public_chat",
				Visibility: "public",
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: initState,
			}

			crr, err := matrix.CreateRoom(req)

			if err != nil || crr == nil {
				log.Println(err)
				break
			}

			log.Println("adding topic")
			resp, err := matrix.SendMessageEvent(resp.Streams.Topics, "m.room.message", map[string]string{
				"body":           topic,
				"formatted_body": topic,
				"title":          topic,
				"msgtype":        "m.text",
				"topic_room_id":  crr.RoomID,
			})
			if err != nil {
				log.Println(err)
			} else {
				log.Println(resp)
			}

		}

	}
}
