import { writable } from 'svelte/store';
import { navigate } from 'svelte-navigator'

function createApp() {

  let app = {
    loadingMessage: 'Calibrating',
    active: false,
    ready: false,
    accounts: [],
    active_account: null,
    settings: {
      switcher: {
        mode: 'normal',
      },
      appearance: {
        theme: window.settings?.appearance?.theme || "dark",
        displayMode: window.settings?.appearance?.displayMode || "cozy",
        fontScaling: window.settings?.appearance?.fontScaling || 16,
        eventSpacing: window.settings?.appearance?.eventSpacing || 16,
      },
      accessibility: {
        saturation: window.settings?.accessibility?.saturation || 100,
        limitDisplayWidth: window.settings?.accessibility?.limitDisplayWidth || false,
        showMessageButton: window.settings?.accessibility?.showMessageButton || false,
      },
      language: window.settings?.language || "en_US",
    },
    events: [],
    temp_events: [],
    members: {},
    allMembers: {},
    allRooms: [],
    allReplies: {},
    isMobile: false,
    mobileViewToggled: true,
    emoji: {
      active: false,
      event: null,
    },
    focusEditorInRoomID: null,
    splitView: false,
    messages: {
      mentionUser: null,
    },
    alerts: [],
    locale: null,
  }

  if(!window.settings) {
    localStorage.setItem('settings', JSON.stringify(app.settings))
  }

  import(`../settings/components/language/locale/${app.settings.language}.js`).then(d => {
    if(d) {
      app.locale = d.default
    }
  })


  let buildServer = (rooms, events, serverID) => {
    let server = {}
    server.rooms = []
    server.members = {}
    events.forEach(event => {
      if(event.type == 'm.space.parent') {
        server.child = true
        server.server_id = event.state_key
      }
      if(event.type == 'm.room.name') {
        server.name = event.content.name
      }
      if(event.type == 'm.room.avatar') {
        server.avatar = event.content.url
      }
      if(event.type == 'm.room.topic') {
        server.topic = event.content.topic
      }
      if(event.type == 'm.room.member') {
        /*
        if(event?.unsigned?.prev_content?.is_direct) {
            server.is_direct = true
          server.room_type = 'dm'
        }
        */
        let user_id;
        if(event?.content?.membership == "join") {
          user_id = event.sender
        } else if(event?.content?.membership =="invite") {
          user_id = event.state_key
        } else if(event?.content?.membership =="leave") {
          user_id = event.state_key
        }



        server.members[user_id] = {
          username: strip(user_id),
          display_name: event?.content?.displayname,
          avatar_url: event?.content?.avatar_url,
        }
          app.allMembers[user_id] = {
            username: strip(user_id),
            display_name: event?.content?.displayname,
            avatar_url: event?.content?.avatar_url,
          }
      }

      if(event.type == 'm.room.create') {
        server.sender = event.sender
        server.origin_server_ts = event.origin_server_ts
        server.created_on = event.unsigned?.age
        if(identity?.user_id == event.sender) {
          server.owner = true
        }
      }
      if(event.type == 'commune.room') {
        server.room_type = event.content.room_type
      }
      if(event.type == 'commune.room.short_alias') {
        server.short_alias = event.content.short_alias
      }
      if(event.type == 'commune.room.thread') {
        server.room_type = event.content.room_type
        server.thread_in_room_id = event.content.thread_in_room_id
        server.expire_thread = event.content.expire_thread
      }
      if(event.type == 'm.room.canonical_alias') {
        let alias = event.content.alias.substring(1)
        alias = alias.split(":")[0]
        server.alias = alias
        alias =`/${alias}`
        server.pathname = alias
      }
      if(event.type == 'm.room.type') {
        server.alt_alias = event.content.alias
      }
      if(event.type == 'm.room.create') {
        if(event.content.creator == identity.user_id) {
          server.owner = true
        }
      }
      if(event.type == 'm.space.child') {
        /*
        let roomID = event.state_key
        for (const [id, room] of Object.entries(rooms)) {
          if(id == roomID) {
            let child = buildServer(rooms, room.state.events)
            child.room_id = id
            child.server_id = serverID
            child.server_pathname = server.pathname
            child.events = []
            child.members = []
            child.pinged = false
            //child.pathname = server.pathname + child.pathname
            server.rooms.push(child)
          }
        }
          server.rooms.push({
            room_id: 
          })
        */

          let room = {
            //room_id: event.state_key,
            channel_id: event.state_key,
            room_id: event?.content?.streams?.[event?.content?.default_stream],
            alias: event.content.local_part,
            pathname: `/${event.content.local_part}`,
            name: event.content.name,
            server_id: serverID,
            sender: event.sender,
            room_type: event.content.default_stream,
            streams: event.content.streams,
          }

          let stored = localStorage.getItem('streams')
          if(stored) {
            let streams = JSON.parse(stored)
            if(streams && (room.channel_id in streams)) {
              room.room_type = streams[room.channel_id]
              room.room_id = room.streams[streams[room.channel_id]]
            }
          }

          if(event.content?.thread) {
            room.thread_in_room_id = event.content.thread.thread_in_room_id
            room.expire_thread = event.content.thread.expire_thread
            room.room_id = event?.state_key
            room.streams = {
              "chat": event?.state_key,
            }
          }


        if(event.content?.default) {
          room.default = true
        }


        server.rooms.push(room)
      }
      if(server.rooms.length > 0) {
        server.rooms?.sort((a, b) => (a.pathname > b.pathname) ? 1 : -1)
        let def = server.rooms.filter(child => child.name == 'general')[0]?.pathname
        server.default_room = server.pathname + def
        server.active_room = server.pathname + def
      }
    })
    return server
  }

  let addServer = (server) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      account.servers.push(server)

      p.allRooms.push(server)

      server.rooms.forEach(room => {
        p.allRooms.push(room)
      })

      navigate(server.default_room)

      return p
    })
  }

  let addRoom = (server, newRoom, dontNavigate, dm) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let xx = account.servers.filter(x => x.pathname == server)[0]
      if(xx) {
        xx.rooms.push(newRoom)
      }
      p.allRooms.push(newRoom)

      if(dm) {
        let xx = account.direct_messages.filter(x => x.alias == server)[0]
        if(xx) {
          xx.rooms.push(newRoom)
          account.direct_messages.push(newRoom)
        }
      }


      if(!dm) {
        if(!dontNavigate) {
          navigate(`${server}${newRoom.pathname}`)
        }
      } else {
        if(!dontNavigate) {
          navigate(`/messages/${newRoom.alias}`)
        }
      }


      return p
    })
  }

  let updateActiveRoom = (server, path) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let xx = account.servers.filter(x => x.pathname == server)[0]
      if(xx) {
        xx.active_room = path
      }



      return p
    })
  }

  let updateActiveHomePage = (path) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      account.home.active_page = path


      return p
    })
  }

  let updateActiveDirectMessages = (path) => {
    update(p => {
      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]

        let ind = account.active_direct_messages?.findIndex(x => x.pathname == path)
        if(ind == -1) {
          account.active_direct_messages.push({
            pathname: path,
            visible: false,
          })
        }


      return p
    })
  }


  let updateActiveRooms = (server, path) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]

      account.active_rooms.forEach(page => {
        page.visible = false
      })
        let ind = account.active_rooms?.findIndex(x => x.pathname == path)
        if(ind == -1) {
          account.active_rooms.push({
            pathname: path,
            visible: false,
          })
          account.active_rooms = account.active_rooms
        }


      return p
    })
  }

  let markRoomVisible = (server, path) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
      account.active_rooms.forEach(page => {
        page.visible = false
      })
        let ind = account.active_rooms?.findIndex(x => x.pathname == path)
        if(ind == -1) {
          account.active_rooms.push({
            pathname: path,
            //visible: true,
          })
        } else {
          //p.active_rooms[ind].visible = true
        }


      return p
    })
  }


  let updateEvents = (props) => {
    let id = props.room_id
    let events = props.events
    let start = props.start
    let end = props.end
    let backfill = props.backfill
    let forward = props.forward
    let init = props.init

    let roomType = props.roomType


    update(p => {






      let ev = []
      if(!backfill) {
        if(events) {
          for (let i = events.length - 1; i >= 0; i--) {
            let e = events[i];
            e = processEvent(e, {
              delivered: true,
            })
            if(roomType == 'topics' || init || forward) {
              ev.unshift(e)
            } else {
              ev.push(e)
            }
          }
        }
      } else {
        events.forEach(e => {
          e = processEvent(e, {
            delivered: true,
          })
            ev.push(e)
        })
      }

      // subsequent updates
      if(backfill) {

        /*
        if(topics) {
          let evv = p.events?.[id]?.topics
          if(evv?.events) {
            ev.forEach(eve => {
              let ind = evv.topics.events.findIndex(x => x.event_id == eve.event_id)
              if(ind == -1) {
                p.events[id].topics.events.unshift(eve)
              }
            })
          }
        } else {
          let evv = p.events?.[id]?.chat
          if(evv?.events) {
            ev.forEach(eve => {
              let ind = evv.chat.events.findIndex(x => x.event_id == eve.event_id)
              if(ind == -1) {
                p.events[id].chat.events.unshift(eve)
              }
            })
          }
        }
        */

          let evv = p.events?.[id]
          if(evv?.events) {
            ev.forEach(eve => {
              let ind = evv.events.findIndex(x => x.event_id == eve.event_id)
              if(ind == -1) {
                if(roomType == 'topics') {
                  p.events[id].events.push(eve)
                } else {
                  p.events[id].events.unshift(eve)
                }
              }
            })
            p.events[id].start = start
            p.events[id].end = end
          }

      } else if(forward) {
          let evv = p.events?.[id]
          if(evv?.events) {
            ev.forEach(eve => {
              let ind = evv.events.findIndex(x => x.event_id == eve.event_id)
              if(ind == -1) {
                p.events[id].events.push(eve)
              }
            })
            p.events[id].start = start
            p.events[id].end = end
          }
      } else {

        if(!(id in p.events)) {
          p.events[id] = {}
        }

        p.events[id] = {
          events: ev,
          start: start,
          end: end,
        }

      }



      window.events = p.events

      return p
    })
  }

  let processMember = (member) => {
    return member
  }

  let updateMembers = (roomID, member) => {
    update(p => {


        let x = {
          username: member.username,
          avatar_url: member.avatar_url,
          display_name: member.display_name,
        }

      if(roomID in p.members) {
        p.members[roomID][member.user_id] = x
      }
      if(!(member.user_id in app.allMembers)) {
        app.allMembers[member.user_id] = {
          username: member.username,
          avatar_url: member.avatar_url,
          display_name: member.display_name,
        }
      }

      return p
    })
  }


  window.printMatrix = () => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
  }



  const { subscribe, set, update } = writable(app);



  let buildAccount = (opts) => {
    update(p => {

      let account = {
        user_id: opts?.identity?.user_id,
        access_token: opts?.identity?.access_token,
        matrix_access_token: opts?.identity?.matrix_access_token,
        home_server: opts?.identity?.home_server,
        username: strip(opts?.identity.user_id),
        display_name: opts?.identity.display_name,
        avatar_url: opts?.identity?.avatar_url,
        status: opts?.identity?.user_status,
        servers: [],
        presence: 'online',
        home: {
          default_page: '/',
          active_page: '/',
        },
        direct_messages: [],
        active_direct_messages: [],
        dm_notifications: [],
        dm_requests: [],
        mentions: [],
        sync_state: opts?.identity?.sync_state,
        notifications: opts?.identity?.notificiations,
        account_data: opts?.identity?.account_data || {},
        active_rooms: [],
      }

        app.allMembers[opts?.identity.user_id] = {
          username: strip(opts?.identity.user_id),
          display_name: opts?.identity.display_name,
          avatar_url: opts?.identity.avatar_url,
        }


      let x = opts?.identity?.user_id.split(":")[0]
      account.username = x?.substring(1)


      //lets join invited rooms first


      if(!opts?.identity?.sync_state) {
        return
      }

      let rooms = opts?.identity?.sync_state?.rooms?.join
      if(rooms) {

        for (const [id, room] of Object.entries(rooms)) {
          let server = buildServer(rooms, room.state.events, id)
          server.room_id = id
          server.channel_id = id
          if(!server?.child && server.name) {
            p.members[id] = server.members
            if(server.pathname && server.room_id) {
              account.servers.push(server)
            }
          }

          if(!server.name) {
            let n = []
            for (const [id, user] of Object.entries(server?.members)) {
              if(id != opts?.identity?.user_id) {
                if(user.display_name?.length > 0) {
                    n.push(user.display_name)
                } else {
                    n.push(strip(id))
                }
              }
            }
            server.name = n.join(', ')
          }
          /*
          if(server.is_direct) {
            account.direct_messages.push(server)
          }
          */
          server.rooms.forEach(x => {
            x.server_alias = server.alias
          })

          //DM rooms here
          if(opts?.identity?.account_data) {
            for (const [_, room] of Object.entries(opts?.identity?.account_data)) {
              let roomID = room[0]
              if(roomID == id) {


                server.room_type = "dm"
                server.alias = server.alt_alias
                server.alt_alias = null

                if(Object.entries(server.members)?.length == 2) {
                  for (const [id, user] of Object.entries(server?.members)) {
                    if(id != account?.identity?.user_id) {
                      server.dm_with = id
                      if(user.avatar_url?.length > 0) {
                        server.avatar_url =  user.avatar_url
                      }
                    }
                  }
                }

                let ind = account.direct_messages.findIndex(x => x.room_id == room)
                if(ind == -1) {
                  account.direct_messages.push(server)
                }

                if(server.rooms?.length > 0) {
                  server.rooms.forEach(x => {
                    x.child = true
                    account.direct_messages.push(x)
                  })
                }


                // initialize DM notificiations in switcher
                for (const [_, item] of Object.entries(opts?.identity?.notifications?.notifications)) {
                  if(item.room_id == id) {
                    let count = 0;
                    if(id in account.dm_notifications) {
                      count = account.dm_notifications[id]?.count
                    }
                    account.dm_notifications[id] = {
                      room_id: id,
                      count: count + 1,
                      origin_server_ts: item?.event?.content?.origin_server_ts,
                    }
                  }
                }

              }
            }


          }

          let stored = localStorage.getItem('streams')
          let streams;
          if(stored) {
            streams = JSON.parse(stored)
          }

          let indd = p.allRooms.findIndex(x => x.room_id == server.room_id)
          if(indd == -1) {
            if(streams && (server.room_id in streams)) {
              server.room_type = streams[server.room_id]
            }
            p.allRooms.push(server)
          }



        }

        account.direct_messages.forEach(dm => {
          p.allRooms.forEach(room => {
            if(room.thread_in_room_id == dm.room_id) {
              let ind = dm.rooms.findIndex(x => x.room_id == room.room_id)
              if(ind == -1) {
                dm.rooms.push(room)
              }
              let indd = account.direct_messages.findIndex(x => x.room_id == room.room_id)
              if(indd == -1) {
                account.direct_messages.push(room)
              }
            }
          })
        })

        account.servers?.sort((a, b) => (a.pathname > b.pathname) ? 1 : -1)
      }

      // Build DM requests
      let invites = opts?.identity?.sync_state?.rooms?.invite
      if(invites) {
        for (const [room_id, invite] of Object.entries(invites)) {

          let state = invite?.['invite_state']

          let isDM = state?.Events?.filter(x => x.state_key == account?.user_id &&
          x.content?.is_direct &&
          x.content?.membership == 'invite')[0]

          if(isDM) {
            let sender = state?.Events?.filter(x => x.type == 'm.room.create')[0]?.sender
            let member = state?.Events?.filter(x => x.type == 'm.room.member' &&
            x.state_key == sender)[0]
            account.dm_requests.push({
              room_id: room_id,
              sender: {
                user_id: sender,
                username: strip(sender),
                display_name: member.content.displayname,
                avatar_url: member.content.avatar_url,
              }
            })
          }
        }
      }





        account.matrix = window.matrixcs.createClient({
          baseUrl: homeServer,
          accessToken: opts?.identity?.matrix_access_token,
          userId: opts?.identity.user_id,
          sessionStore: new window.matrixcs.WebStorageSessionStore(localStorage),
          cryptoStore: new window.matrixcs.IndexedDBCryptoStore(indexedDB, 'crypto-store'),
          deviceId: opts?.identity?.device_id,
        });



        account.matrix.on("event", function(e) {
          if(e.getType() == 'm.room.message' || 
            e.getType() == 'm.reaction') {
            const roomID = e.event.room_id
            if(e?.event?.transaction_id?.length > 0) {
              p.ready = true
            }

            if(account.user_id == p.active_account) {


              if(e.event.content?.msgtype == 'commune.room.thread.new') {
                if(p.ready) {
                  p.temp_events.push(e.event)
                }
              } else {
                if(p.ready) {
                  addEventToRoom(roomID, e.event)
                }
              }
            }
          }
          if(e.getType() == 'm.room.create' ) {
            if(p.ready) {
              if(e?.event?.sender != p.active_account) {
                console.log(e.event)

                /*
                let sameHS = eventFromHomeServer(e.event?.room_id)

                // Check if this is a DM invite
                if(sameHS) {

                  getRoomState(e?.event?.room_id).then(res => {
                    let invited = res?.filter(x => x.type == 'm.room.member' &&
                    x?.state_key == opts?.identity?.user_id &&
                    x?.content?.is_direct &&
                    x?.content?.membership == 'invite')[0]

                    if(invited) {

                      let sender = res?.filter(x => x.type == 'm.room.create')[0]?.sender
                      let member = res?.filter(x => x.type == 'm.room.member' &&
                      x.state_key == sender)[0]
                        console.log(member)
                      addDMRequest({
                        room_id: e.event.room_id,
                        sender: {
                          user_id: sender,
                          username: strip(sender),
                          display_name: member.content.displayname,
                          avatar_url: member.content.avatar_url,
                        }
                      })

                    }
                  })

                }
                */
                /*
                getRoomState(e?.event?.room_id).then(res => {
                  console.log(res)
                  let invited = res?.filter(x => x.type == 'm.room.member' &&
                  x?.state_key == opts?.identity?.user_id &&
                  x?.content?.is_direct &&
                  x?.content?.membership == 'invite')[0]
                  return [res, invited]

                }).then(([res, invited]) => {
                  console.log(res)
                  console.log(invited)
                  console.log(invited)
                  console.log(invited)
                })
                  joinRoom({
                    room_id: e.event.room_id,
                    token: account?.matrix_access_token,
                    navigate: false,
                  })
                //*/

              }
            }
          }

          if(e.getType() == 'm.room.member' ) {
            const roomID = e.event.room_id
            if(p.ready) {

              //Check if this is DM request

              let invited = e.event.type == 'm.room.member' &&
              e.event?.state_key == opts?.identity?.user_id &&
              e.event?.content?.is_direct &&
              e.event?.content?.membership == 'invite'

              if(invited) {

                let sender = e.event.sender

                fetchUserProfile(sender).then(res => {
                  return res
                }).then(res => {
                  addDMRequest({
                    room_id: e.event.room_id,
                    sender: {
                      user_id: sender,
                      username: strip(sender),
                      display_name: res.displayname,
                      avatar_url: res.avatar_url,
                    }
                  })
                })


              }




              if(e.event?.content?.membership == 'join') {

                let mem = {
                    user_id: e.event.sender,
                    username: strip(e.event.sender),
                    avatar_url: e.event?.content?.avatar_url,
                    display_name: e.event?.content?.display_name,
                }

                //addEventToRoom(e.event?.room_id, e.event)
                updateMembers(e.event?.room_id, mem)
                app.allMembers[e.event.sender] = {
                  username: strip(e.event.sender),
                  avatar_url: e.event?.content?.avatar_url,
                  display_name: e.event?.content?.display_name,
                }
              }
            }
          }

          if(e.getType() == 'commune.room.thread' ) {
          }
          if(e.getType() == 'm.room.redaction' ) {
            const roomID = e.event.room_id
              redactEvent(roomID, e.event)
          }

        });

        account.matrix.on("sync", function(e) {
          p.ready = true 
        })



        account.matrix.startClient({
          initialSyncLimit: 0,
        })

      p.accounts.push(account)


      /*
      let invite = opts?.identity?.sync_state?.rooms?.invite
      if(invite) {
        for (const [roomID, _] of Object.entries(invite)) {
            joinRoom({
              room_id: roomID,
              token: opts?.identity?.matrix_access_token,
              navigate: false,
            })
        }
      }
      */

      return p
    })

  }

  let addDMRequest = (request) => {
    update(p => {
      console.log(request)
      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]
      account.dm_requests.push(request)
      return p
    })
  }

  let rejectDMRequest = (room_id) => {
    update(p => {
      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let ind = account.dm_requests.findIndex(x => x.room_id == room_id)
      if(ind != -1) {
        account.dm_requests?.splice(ind, 1)
        leaveRoom(room_id).then(res => {
          console.log(res)
        })
      }
      return p
    })
  }

  let acceptDMRequest = (room_id) => {
    update(p => {
      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let ind = account.dm_requests.findIndex(x => x.room_id == room_id)
      if(ind != -1) {
        account.dm_requests?.splice(ind, 1)
        joinRoom({
          room_id: room_id,
          token: account?.matrix_access_token,
          navigate: true,
        })
      }
      return p
    })
  }


async function leaveRoom(room_id) {

    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${room_id}/leave`

    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]

    let resp = await fetch(endpoint, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


  let activate = () => {
    update(p => {
      let mode = localStorage.getItem("switcher-mode");
      if(mode) {
        if(mode == 'normal' || mode == 'expanded' || mode == 'collapsed') {
          p.settings.switcher.mode = mode
        }
      }

      let isMobile = navigator?.userAgentData?.mobile
      let vw = window.innerWidth <= 600

      let userAgentMobile = () => {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
      }

      if(isMobile || vw || userAgentMobile()) {
        p.isMobile = true
      }


      buildAccount({
        identity: identity,
      })


      if(identity?.alt_accounts) {
        identity.alt_accounts?.forEach(account => {
          buildAccount({
            identity: account,
          })
        })
      }

      /*
      let active_account = localStorage.getItem('active_account')
      if(active_account) {
        let account = p?.accounts?.filter(a => a.user_id == active_account)[0]
        if(account) {
          p.active_account = account.user_id
        }
      } else {
        p.active_account = identity?.user_id
      }
      */
      p.active_account = identity?.user_id

      p.active = true


      window.printStore = () => {
        console.log(app)
      }

      return p
    })
  }


  let addNewAccount = (opts) => {
    update(p => {

      buildAccount({
        identity: opts.identity,
      })


      return p
    })
  }

  let switchToAccount = (opts) => {
    update(p => {

      p.active_account = opts.user_id

      navigate(`/`)

      localStorage.setItem("active_account", opts?.user_id);


      return p
    })
  }



  let initSocket = (account, endpoint, wait) => {
    setTimeout(() => {
      account.socket = new WebSocket(endpoint)
        account.socket.onopen = () => {
            console.log("Sync socket opened.")
        };
        account.socket.onclose = () => {
          console.error("yikes, sync socket closed. reconnecting in a bit...")
          initSocket(account, endpoint, 5000 + wait)
        };
        account.socket.onmessage = (e) => {
          let data = JSON.parse(e.data)
          console.log("lol we got new data", data)
        }
        account.socket.onerror = (e) => {
          console.error("yikes, sync socket error", e)
          initSocket(account, endpoint, 5000 + wait)
        }

    }, wait)
  }



async function initCrypto(matrix) {

        await matrix.initCrypto((err, res) => {
          console.log(res)
        })
}


  let processEvent = (event, props) => {
    if(props?.delivered) {
      event.delivered = true
    } else {
      event.delivered = false
    }
    return event
  }

  let addEventToRoom = (room_id, event, topics) => {
    update(p => {
      let opts = {delivered: true}
      if(event.transaction_id) {
        opts.delivered = false
      }
      event = processEvent(event, opts)



      let roomType;

      if(event?.content?.topic) {
        roomType = 'topics'
      } else {
        roomType = 'chat'
      }


      if(roomType == 'topics') {
        event.fresh = true
      }

      if(room_id in p.allReplies) {
        if(event.content['m.new_content']) {
          let event_id = event?.content?.[`m.relates_to`]?.event_id

          let index = p.allReplies[room_id]?.findIndex(x => x.event_id == event_id)
          if(index != -1) {
            p.allReplies[room_id][index].content.title = event.content?.[`m.new_content`]?.title
            p.allReplies[room_id][index].content.body = event.content?.[`m.new_content`]?.body
            p.allReplies[room_id][index].content.formatted_body = event.content?.[`m.new_content`]?.formatted_body
            p.allReplies[room_id][index].edited = true
          }
        }
      }

      if((room_id in p.events)) {


          let ind = p.events[room_id].events.findIndex(x => x.transaction_id == event.unsigned?.transaction_id)
          if(ind == -1 || !event?.unsigned?.transaction_id) {

            // This is a new message/post.
            if(!event.content['m.new_content']) {
              if(roomType == 'topics') {
                p.events[room_id].events.unshift(event)
              } else {
                p.events[room_id].events.push(event)
              }
            // This is an edited message/post.
            } else {
              let event_id = event?.content?.[`m.relates_to`]?.event_id

              let index = p.events[room_id]?.events?.findIndex(x => x.event_id == event_id)
              if(index != -1) {
                p.events[room_id].events[index].content.title = event.content?.[`m.new_content`]?.title
                p.events[room_id].events[index].content.body = event.content?.[`m.new_content`]?.body
                p.events[room_id].events[index].content.formatted_body = event.content?.[`m.new_content`]?.formatted_body
                p.events[room_id].events[index].edited = true
              }



            }

          } else {
            p.events[room_id].events[ind].age = event.age
            p.events[room_id].events[ind].unsigned.age = event.unsigned.age
            p.events[room_id].events[ind].event_id = event.event_id
            p.events[room_id].events[ind].origin_server_ts = event.origin_server_ts
            p.events[room_id].events[ind].delivered = true
          }



        /*
      } else {
        p.events[room_id] = {events:[]}
        if(!event.content['m.new_content']) {
          if(topics){
            p.events[room_id].events.unshift(event)
          } else {
            p.events[room_id].events.push(event)
          }
        }
        */
      }

      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]

      let ind = p.allRooms.findIndex(x => x.room_id == room_id)

      if(ind != -1 ) {
        let rm = p.allRooms[ind]
        let isDM = rm?.room_type == "dm"


        if(isDM && event.sender != p.active_account) {
          let alias = `/messages/${rm.alias}`
          let isActive = window.location.pathname == alias

          if(!isActive) {
            let count = 0;
            if(room_id in account.dm_notifications) {
              count = account.dm_notifications[room_id]?.count
            }
            account.dm_notifications[room_id] = {
              room_id: room_id,
              count: count + 1,
              origin_server_ts: event?.origin_server_ts,
            }
            let sound = document.createElement("audio");
            sound.src = `/static/audio/bell.ogg`;
            sound.setAttribute("preload", "auto");
            sound.setAttribute("controls", "none");
            sound.style.display = "none";
            document.body.appendChild(sound);
            sound.play();
          }
        }
      }

      let div = document.createElement('div')
      div.innerHTML = event.content.formatted_body
      let spans = div.querySelectorAll('span')
      if(spans?.length > 0) {
          spans.forEach(span => {
              if(span.dataset.userid == p.active_account) {
                console.log("WE SHOULD ALERT user")
              }
          })
      }

      //add topic replies
      let isReply = event?.content?.topic_reply

      if(isReply && (event?.room_id in p.allReplies)) {
        console.log("heheh")
        event.fresh = true
        let ind = p.allReplies[event?.room_id]?.findIndex(x => x.event_id == event.event_id)
        if(ind == -1) {
          p.allReplies[event?.room_id].unshift(event)
        }
      }

      //add topic reactions
      let isReaction = event?.type == `m.reaction`

      if(isReaction && (event?.room_id in p.allReplies)) {
        event.fresh = true
        let ind = p.allReplies[event?.room_id]?.findIndex(x => x.event_id == event.event_id)
        if(ind == -1) {
          p.allReplies[event?.room_id].unshift(event)
        }
      }

      return p
    })
  }

  let resetDMnotification = (id) => {
    update(p => {

      let account = p?.accounts?.filter(account => account.user_id == p.active_account)[0]

      if(id in account.dm_notifications) {
        delete account.dm_notifications[id]
      }
      return p
    })
  }


  let redactEvent = (room_id, event) => {
    update(p => {
      console.log("redacting event", event)
      if(room_id in p.events) {
        let ind = p.events[room_id]?.['chat']?.events.findIndex(x => x.event_id == event.redacts)
        if(ind != -1) {
          p.events[room_id]['chat'].events[ind].redacted = true
          p.events[room_id]['chat'].events[ind].redacted_because = event?.content?.reason
          p.events[room_id]['chat'].events[ind].content = {}
        }
      }
      return p
    })
  }

  let removeEvent = (room_id, event) => {
    update(p => {

      if(event?.room_id in p.allReplies) {

        let ind = p.allReplies[event?.room_id].findIndex(x => x.event_id == event.event_id)
        console.log(ind)
        if(ind != -1) {
          p.allReplies[event?.room_id]?.splice(ind, 1)
        }
      }

      if(event?.room_id in p.events) {

        let ind = p.events[event?.room_id]?.events.findIndex(x => x.event_id == event.event_id)
        console.log(ind)
        if(ind != -1) {
          p.events[event?.room_id]?.events?.splice(ind, 1)
        }
      }

      return p
    })
  }


  let updateReactions = (roomID, event) => {
    update(p => {
      if(roomID in p.events) {
        const evid = event?.content?.['m.relates_to']?.event_id
        let ind = p.events[roomID].events.findIndex(x => x.event_id == evid)
        if(ind != -1) {
          let e = p.events[roomID].events[ind]?.unsigned?.['m.relations']?.['m.annotation']
          if(e?.chunk) {
            let ind = e.chunk.findIndex(x => x.key == event?.content?.['m.relates_to'].key)
            console.log(ind)
            if(ind == -1) {
              e.chunk.push({
                count: 1,
                key: event?.content?.['m.relates_to'].key,
                type: "m.reaction",
              })
            } else {
              e.chunk[ind].count = e.chunk[ind]?.count + 1
            }
          } else {
            p.events[roomID].events[ind].unsigned['m.relations'] = {
              'm.annotation': {
                chunk: [
                  {
                    count: 1,
                    key: event?.content?.['m.relates_to'].key,
                    type: "m.reaction",
                  }
                ]
              }
            }
          }
        }
      }
      return p
    })
  }


  let loadIdleRoomsEvents = (server) => {
      update(p => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
      let rooms = account.servers.filter(x => x.pathname == server)[0]?.rooms
      console.log(rooms)
      if(rooms) {
        rooms.forEach(room => {
          let topics = room?.room_type == 'topics'
          fetchRoomEvents({room_id: room?.room_id}).then(res => {
              let props = {
                  room_id: room?.room_id,
                  events: res.chunk,
                  start: res.start,
                  end: res.end,
                  backfill: false,
                  isTopics: topics,
              }
              store.updateEvents(props)
          })
          /*
          getMembers(room?.room_id).then(res => {
              store.updateMembers(room?.room_id, res.joined)
          })
          */
        })
        Object.keys(rooms).forEach((roomID) => {
          /*
          fetchRoomEvents(roomID).then(res => {
            console.log(res)
              store.updateEvents(roomID, res.chunk, res.start, res.end, false)
          })
          */
        });
      }


      return p
    })
  }



  let syncTemporaryEvent = (id, event) => {
    update(p => {
      if(id in p.events) {
        let ind = p.events[id].events.findIndex(x => x.client_ts == event.content.client_ts)
        if(ind != -1) {
          p.events[id].events[ind].age = event.age
          p.events[id].events[ind].unsigned.age = event.unsigned.age
          p.events[id].events[ind].event_id = event.event_id
          p.events[id].events[ind].origin_server_ts = event.origin_server_ts
          p.events[id].events[ind].delivered = true
        }
      }
      return p
    })
  }

async function fetchRoomEvents(opts) {
  let roomID = opts.room_id
  let start = opts.start
  let roomType = opts.roomType

  let optDir = opts.dir


  let msgType = `m.room.message`
  let reactionType = `m.reaction`

  if(roomType == 'topics') {
    msgType = `commune.room.topics.post`
    reactionType = `commune.room.topics.post.reaction`
  }

  if(opts.topicPost) {
    msgType = opts.msgType
    reactionType = opts.reactionType
  }


  let filterContent = {types: [msgType, reactionType]}

  let filter = JSON.stringify(filterContent)

  let limit = 50

  /*
  if(roomType == 'topics') {
    limit = 7
  }

    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/messages?limit=${limit}&dir=b&filter=${filter}`
  if(start) {
    endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/messages?limit=${limit}&dir=b&from=${start}&filter=${filter}`
  }
  */

  let dir = `b`

  if(optDir == `f`) {
    dir = `f`
  }

    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/messages?limit=${limit}&dir=${dir}`
  if(start) {
    endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/messages?limit=${limit}&dir=${dir}&from=${start}`
  }
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]

    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

  let fetchRoomMembers = (roomID) => {
    update(p => {
      fetchMembers(roomID).then(res => {
        console.log(res)
      })
      return p
    })
  }

async function fetchEventContext(opts) {
  let roomID = opts.room_id
  let eventID = opts.event_id


    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/context/${eventID}?limit=50`

    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]

    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function createTopicRoom(roomID) {
    let endpoint = `/topic/create`
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
      method: 'POST',
        headers: {
            'Authorization': account.access_token
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


async function fetchMembers(roomID) {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/joined_members`
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function fetchNotifications() {
    let endpoint = `${homeServer}/_matrix/client/r0/notifications`
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function updateReadReceipts(roomID, eventID, type) {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/receipt/${type}/${eventID}`
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function updateReadMarker(roomID, eventID) {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/read_markers`

  let data = {
    ['m.fully_read']: eventID,
    ['m.read']: eventID,
  }
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
        method: 'POST',
      body: JSON.stringify(data),
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function joinRoom(opts) {

  const roomID = opts.room_id
  const token = opts.token
  const navigateTo = opts.navigate
  const federated = opts.federated

  userJoinRoom(roomID, token, federated).then(res => {
    if(res) {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
      console.log(res)
      getRoomState(roomID).then(res => {
        if(res) {
          console.log(res)
          let isDirect = false;
          let creator;
          res?.forEach(event => {
            if(event.type == "m.room.create") {
              creator = event.content.creator
            }
          })
          console.log("checking for state events")
          console.log("checking for state events")
          console.log("checking for state events")
          res?.forEach(event => {
            console.log(event)
            if(event.type == 'm.room.member' && event.sender == app.active_account) {
              if(event?.unsigned?.prev_content?.is_direct ||
              event?.content?.membership == 'join') {
                isDirect = true
                addDirectAccountData(roomID, creator, account.account_data)
              }
            }
          })

          let room = buildServer({
            [roomID]: res,
          }, res)
          room.room_id = roomID
          if(isDirect) {
            room.is_direct = true
            room.room_type = 'dm'
          }
          let server = account?.servers?.filter(server => server.room_id == room.server_id)[0]
          room.server_pathname = server?.pathname
          if(server && !isDirect) {
            let ind = server?.rooms?.findIndex(x => x.room_id == roomID)
            if(ind == -1) {
              room.server_alias = server.alias
              server?.rooms.push(room)
              let indd = app.allRooms.findIndex(x => x.room_id == roomID)
              if(indd == -1) {
                app.allRooms.push(room)
              }
            }
          }

          if(isDirect) {

            if(Object.entries(room.members)?.length == 2) {
              for (const [id, user] of Object.entries(room?.members)) {
                if(id != app.active_account) {
                  room.dm_with = id
                  if(user.avatar_url?.length > 0) {
                    room.avatar_url =  user.avatar_url
                  }
                }
              }
            }

            if(!room.name) {
              let n = []
              for (const [id, user] of Object.entries(room.members)) {
                if(id != app.active_account) {
                  if(user.display_name?.length > 0) {
                      n.push(user.display_name)
                  } else {
                      n.push(strip(id))
                  }
                }
              }
              room.name = n.join(', ')
            }

            room.alias = room.alt_alias
            room.new = true
            let ind = account?.direct_messages?.findIndex(x => x.room_id == roomID)
            if(ind == -1) {
              account?.direct_messages.push(room)
              let indd = app.allRooms.findIndex(x => x.room_id == room.room_id)
              if(indd == -1) {
                app.allRooms.push(room)
              }
            }

            //
            if(navigateTo) {
                navigate(`/messages/${room.alias}`)
                updateActiveHomePage(`/messages/${room.alias}`)
                updateActiveDirectMessages(`/messages/${room.alias}`)
            }


          }

          //add new thread to DM
            account.direct_messages.forEach(dm => {
              if(room.thread_in_room_id == dm.room_id) {
                let ind = dm.rooms.findIndex(x => x.room_id == room.room_id)
                if(ind == -1) {
                  dm.rooms.push(room)
                }
                let indd = account.direct_messages.findIndex(x => x.room_id == room.room_id)
                if(indd == -1) {
                  account.direct_messages.push(room)
                }
              }
            })


        }


      })


          fetchRoomEvents({room_id: roomID}).then(res => {
              let props = {
                  room_id: roomID,
                  events: res.chunk,
                  start: res.start,
                  end: res.end,
                  backfill: false,
                  isTopics: false,
              }
              store.updateEvents(props)
          })



    }
  })

}

function addDirectAccountData(roomID, creator, data) {
  data[creator] = [roomID]
  addUserAccountData("m.direct", data).then(res => {
  })
}

async function addUserAccountData(type, data) {
    let account = app?.accounts?.filter(x => x.user_id == app.active_account)[0]
    let endpoint = `${homeServer}/_matrix/client/r0/user/${app.active_account}/account_data/${type}`
    let resp = await fetch(endpoint, {
        method: 'PUT',
        body: JSON.stringify(data),
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function userJoinRoom(roomID ,token, federated) {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/join`
  if(federated) {
    endpoint = `${homeServer}/_matrix/client/r0/join/${roomID}`
  }
    let resp = await fetch(endpoint, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`

        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function getRoomState(roomID) {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/state`
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let resp = await fetch(endpoint, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


  let addNewDM = (roomID, users) => {
    update(p => {


      getRoomState(roomID).then(res => {
        if(res) {

          let room = buildServer({
            [roomID]: res,
          }, res)
          room.room_id = roomID
          room.is_direct = true
          room.room_type = 'dm'
          room.alias = room.alt_alias

          let account = p?.accounts?.filter(account => account.user_id == app.active_account)[0]

          if(!room.name) {
            let n = []
            users.forEach(user => {
                if(user.display_name?.length > 0) {
                    n.push(user.display_name)
                } else {
                    n.push(strip(id))
                }
            })
            room.name = n.join(', ')
          }

          let ind = account?.direct_messages?.findIndex(x => x.room_id == roomID)
          if(ind == -1) {
            account?.direct_messages.push(room)
          }

          fetchRoomEvents({room_id: roomID}).then(res => {
              let props = {
                  room_id: room?.room_id,
                  events: res.chunk,
                  start: res.start,
                  end: res.end,
                  backfill: false,
                  isTopics: false,
              }
              store.updateEvents(props)
          })

          navigate(`/messages/${room.alias}`)
          updateActiveHomePage(`/messages/${room.alias}`)
          updateActiveDirectMessages(`/messages/${room.alias}`)
        }
      })




      return p
    })
  }


async function fetchUserProfile(user_id) {
    let endpoint = `${homeServer}/_matrix/client/r0/profile/${user_id}`
    let resp = await fetch(endpoint)
    const ret = await resp.json()
    return Promise.resolve(ret)
}


  let kill = () => {
    update(p => {
      p.active = false
      return p
    })
  }

  let activeAccount = () => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      return account
    })
  }



  let updateUnreadCount = (id) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let room = server?.rooms.filter(room => room.room_id == id)[0]
        if(room) {
          if(isNaN(room?.unread_count)) {
            room.unread_count = 1
          } else {
            room.unread_count = room.unread_count + 1
          }
        }
      })
      return p
    })
  }

  let updateEventContent = (roomID, eventID, body, unsigned) => {
    update(p => {
      if(roomID in p.events) {
        let ind = p.events[roomID].events.findIndex(e => e.event_id == eventID)
        if(ind != -1) {
          p.events[roomID].events[ind].content.body = body
          p.events[roomID].events[ind].unsigned['m.relations'] = unsigned
        }
      }
      return p
    })
  }




  let toggleSwitcher = () => {
    update(p => {
      let mode = p.settings.switcher.mode
      if(mode == 'normal') {
        p.settings.switcher.mode = 'expanded'
        localStorage.setItem("switcher-mode", 'expanded');
      } else if(mode == 'expanded') {
        p.settings.switcher.mode = 'collapsed'
        localStorage.setItem("switcher-mode", 'collapsed');
      } else {
        p.settings.switcher.mode = 'normal'
        localStorage.setItem("switcher-mode", 'normal');
      }
      p.count += p.count
      return p
    })
  }


let newThread = (room, event, content) => {
  update(p => {

      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room.room_id)[0]
        if(r) {
          r.thread = {
            active: true,
            creating: true,
          }

          if(event) {
            r.thread.events = []
            r.thread.events.push(event)
          } else {
            r.thread.new = true
          }

          if(content) {
            r.thread.content = content
          }

        }
      })
    //this is for DM threads
        let r = account?.direct_messages.filter(r => r.room_id == room.room_id)[0]
        if(r) {
          r.thread = {
            active: true,
            creating: true,
            dm: true,
          }

          if(event) {
            r.thread.event = event
          } else {
            r.thread.new = true
          }

          if(content) {
            r.thread.content = content
          }

        }

    return p
  })
}

let toggleThreadEvent = (room_id, event) => {
  update(p => {

      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room_id)[0]
        if(r?.thread?.active) {
          let ind = r.thread.events?.findIndex(x => x.event_id == event.event_id)
          if(ind == -1) {
            r.thread.events.push(event)
          } else {
            r.thread.events.splice(ind, 1)
          }
          r.thread.events?.sort((a, b) => (a.origin_server_ts > b.origin_server_ts) ? 1 : -1)

        }
      })

    return p
  })
}

let selectThreadEvents = (room_id, events) => {
  update(p => {

      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room_id)[0]
        if(r?.thread?.active) {

          events.forEach(event => {
            let ind = r.thread.events?.findIndex(x => x.event_id == event.event_id)
            if(ind == -1) {
              r.thread.events.push(event)
            }
          })

          r.thread.events?.sort((a, b) => (a.origin_server_ts > b.origin_server_ts) ? 1 : -1)

        }
      })

    return p
  })
}


let openThread = (room, event, threadRoom, isDM) => {
  update(p => {

      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]

    if(isDM) {

      let x = account?.direct_messages.forEach(r => {
        if(r.room_id == room.room_id) {
          r.thread = {
            active: true,
            creating: false,
            viewing: true,
            room: room,
            thread_room: threadRoom,
            event: event,
            new: false,
            content: false,
          }
        }
      })
    } else {



      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room.room_id)[0]
        if(r) {
          r.thread = {
            active: true,
            creating: false,
            viewing: true,
            room: room,
            thread_room: threadRoom,
            event: event,
            new: false,
            content: false,
          }
        }
      })
    }
    return p
  })
}


let discardNewThread = (room, dm) => {
  update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room.room_id)[0]
        if(r) {
          r.thread = {
            active: false,
            creating: false,
            event: null,
            new: false,
          }
        }
      })
    if(dm) {
      let r = account?.direct_messages.filter(r => r.room_id == room.room_id)[0]
      if(r) {
        r.thread = {
          active: false,
          creating: false,
          event: null,
          new: false,
        }
      }
    }
    return p
  })
}

let closeThread = (room) => {
  update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let x = account?.servers.forEach(server => {
        let r = server?.rooms.filter(r => r.room_id == room.room_id)[0]
        if(r) {
          r.thread = {
            active: false,
            viewing: false,
            event: null,
            room: null,
            thread_room: null,
            new: false,
          }
        }
      })
      account?.direct_messages.forEach(r => {
        if(r.room_id == room.room_id) {
          r.thread = {
            active: false,
            viewing: false,
            event: null,
            room: null,
            thread_room: null,
            new: false,
          }
        }
      })
    return p
  })
}

function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}


  let getServer = (pathname) => {
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    return account?.servers?.filter(server => server.pathname === pathname)[0]
  }

  let getRoomID = (server, roomAlias) => {
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let sp =  account?.servers?.filter(x => x.pathname === server)[0]
    return sp?.rooms?.filter(x => x.pathname === roomAlias)[0].room_id
  }


  let getRoom = (server, roomAlias) => {
    let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
    let sp =  account?.servers?.filter(x => x.pathname === server)[0]
    return sp?.rooms?.filter(x => x.pathname === roomAlias)[0]
  }


  let toggleRoomStream = (server_id, room_id) => {
      update(p => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
      let server = account.servers.filter(x => x.room_id == server_id)[0]

      let room = server.rooms.filter(x => x.room_id == room_id)[0]

      let switchTo;

      if(room.room_type == "chat") {
        switchTo = 'topics'
      } else {
        switchTo = 'chat'
      }

      room.room_type = switchTo
      room.room_id = room.streams[switchTo]

      let rm = p.allRooms.filter(x => x.channel_id == room_id)[0]
      if(rm) {
        rm.room_type = switchTo
        room.room_id = room.streams[switchTo]
      }

      let stored = localStorage.getItem('streams')
      if(stored) {
        let streams = JSON.parse(stored)
        streams[room.channel_id] = switchTo
        localStorage.setItem('streams', JSON.stringify(streams))
      } else {
        let streams = {
          [room.channel_id]: switchTo,
        }
        localStorage.setItem('streams', JSON.stringify(streams))
      }


      return p
    })
  }

  let toggleRoomStreamTo = (server_id, room_id, switchTo) => {
      update(p => {
      let account = app?.accounts?.filter(account => account.user_id == app.active_account)[0]
      let server = account.servers.filter(x => x.room_id == server_id)[0]

      let room = server.rooms.filter(x => x.room_id == room_id)[0]

      room.room_type = switchTo
      room.room_id = room.streams[switchTo]

      let rm = p.allRooms.filter(x => x.channel_id == room_id)[0]
      if(rm) {
        rm.room_type = switchTo
        room.room_id = room.streams[switchTo]
      }

      return p
    })
  }


  let updateTopicReplies = (id, events) => {
    update(p => {
      if(!(id in p.allReplies)) {
        p.allReplies[id] = events
      } else {
        events.forEach(event => {
          let ind = p.allReplies[id]?.findIndex(x => x.event_id == event.event_id)
          if(ind == -1) {
            p.allReplies[id].unshift(event)
          }
        })
      }
      return p
    })
  }

  let updateTopicRepliesCount = (server_alias, channel_alias, event_id) => {
    update(p => {
      let account = app?.accounts?.filter(account => account.user_id == p.active_account)[0]
      let server = account?.servers?.filter(x => x.alias == server_alias)[0]
      let room = server?.rooms?.filter(x => x.alias == channel_alias)[0]
      return p
    })
  }


  let mentionUser = (opts) => {
    update(p => {
      p.messages.mentionUser = {
        room_id: opts.room_id,
        username: opts.username,
      }
      return p
    })
  }


  let resetMentionUser = () => {
    update(p => {
      p.messages.mentionUser = null
      return p
    })
  }

  let newAlert = (alert) => {
    update(p => {
      p.alerts.push(alert)
      return p
    })
  }

let homeServerFromUserID = (id) => {
  let items = id?.split(':')
  items = items.slice(1)
  return items.join('')
}

let eventFromHomeServer = (room_id) => {
  return homeServerFromUserID(room_id) == homeServerFromUserID(app.active_account)
}

  let updateTheme = (theme) => {
    update(p => {
      p.settings.appearance.theme = theme
      return p
    })
  }

  let updateDisplayMode = (mode) => {
    update(p => {
      p.settings.appearance.displayMode = mode
      return p
    })
  }
  let updateFontScaling = (scaling) => {
    update(p => {
      p.settings.appearance.fontScaling = scaling
      return p
    })
  }

  let updateEventSpacing = (spacing) => {
    update(p => {
      p.settings.appearance.eventSpacing = spacing
      return p
    })
  }

  let updateLanguage = (lang, data) => {
    update(p => {
      p.settings.language = lang
      p.locale = data

      return p
    })
  }


  let updateSaturation = (val) => {
    update(p => {
      p.settings.accessibility.saturation = val
      saveSettings()
      return p
    })
  }

  let saveSettings = () => {
    update(p => {
      localStorage.setItem('settings', JSON.stringify(p.settings))
      return p
    })
  }




  return {
    subscribe,
    set,
    activate,
    kill,
    toggleSwitcher,
    getServer,
    addServer,
    addRoom, 
    updateActiveRoom,
    activeAccount,
    addEventToRoom,
    updateActiveRooms,
    updateEvents,
    updateMembers,
    getRoom,
    getRoomID,
    fetchRoomEvents,
    updateUnreadCount,
    updateEventContent,
    newThread,
    openThread,
    closeThread,
    discardNewThread,
    updateReactions,
    loadIdleRoomsEvents,
    updateActiveHomePage,
    updateActiveDirectMessages,
    addUserAccountData,
    addNewDM,
    joinRoom,
    resetDMnotification,
    updateReadReceipts,
    updateReadMarker,
    toggleRoomStream,
    toggleRoomStreamTo,
    updateTopicReplies,
    updateTopicRepliesCount,
    removeEvent,
    createTopicRoom,
    toggleThreadEvent,
    selectThreadEvents,
    fetchEventContext,
    mentionUser,
    resetMentionUser,
    newAlert,
    addNewAccount,
    switchToAccount,
    rejectDMRequest,
    acceptDMRequest,
    updateTheme,
    updateDisplayMode,
    updateFontScaling,
    updateEventSpacing,
    updateLanguage,
    updateSaturation,
    saveSettings,
  };
}

export const store = createApp();
