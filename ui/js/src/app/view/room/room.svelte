<script>
import { store } from '../../store/store.js'
import {onMount } from 'svelte'
import { useLocation } from 'svelte-navigator'
import { hourglass } from '../../utils/icons.js'
import Events from './events/events.svelte'
import NewPost from './topics/new-post.svelte'

import LoadingEvents from './placeholder/loading-events.svelte'


export let page;
export let threadView;

const location = useLocation()

let indicate = false;

$: {
    indicate = $location.pathname == page.pathname
}

$: server = page.pathname.split("/")[1]
$: roomAlias = page.pathname.split("/")[2]

$: eventAlias = $location.pathname.split("/")[4]

$: mentionedEvent = messages?.filter(x => x.event_id == eventAlias)

$: mentionedEventExists = eventAlias?.length > 0 && mentionedEvent != null && mentionedEvent != undefined


$: isMessageEvent = $location.pathname.split("/")[3] == `message`

$: isDirectMessage = server == 'messages'

let chat;


let active = false;
let loading = true;

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: servers = account.servers

$: rooms = servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms

$: dms = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.direct_messages

$: dmRoom = dms?.filter(x => x.alias == roomAlias) [0]

$: room = isDirectMessage ? dmRoom : rooms?.filter(r => r.pathname == `/${roomAlias}`)[0]
//$: room = $store.allRooms?.filter(r => r.alias == roomAlias)[0]

$: roomID = room?.room_id


let visible = false;

onMount(() => {

    if(eventAlias && isMessageEvent) {
        let opts = {
            room_id: room.streams?.['chat'],
            event_id: eventAlias,
        }
        store.fetchEventContext(opts).then(res => {
            console.log(res)
            let events = []
            let after = []

            res.events_after?.forEach(event => {
                after.unshift(event)
            })
            after?.forEach(event => {
                events.push(event)
            })
            events.push(res?.event)
            res.events_before?.forEach(event => {
                events.push(event)
            })
            let props = {
                room_id: room.streams?.['chat'],
                events: events,
                start: res.start,
                end: res.end,
                backfill: false,
                roomType: 'chat',
            }
            store.updateEvents(props)
            loading = false
            fetched = true
            visible = true
        })
        return
    }

    if(isMessageEvent && topics) {
        store.toggleRoomStreamTo(room?.server_id, room?.room_id, 'chat')
    }

    let events = $store.events[room?.room_id]
    if(events) {
        loading = false
        fetched = true
        visible = true
        return
    }

    if(!roomID) {
        return
    }

    let roomType = room.room_type


    let opts = {
        room_id: roomID,
        roomType: roomType,
    }
    if(roomType == 'dm') {
        opts.room_id = room.room_id
        opts.roomType = 'chat'
    }

    if(roomType == 'thread') {
        opts.room_id = room.channel_id
        opts.roomType = 'chat'
    }

    store.fetchRoomEvents(opts).then(res => {
        console.log("test")
        console.log(res)
        console.log("test")
        let props = {
            room_id: roomID,
            events: res.chunk,
            start: res.start,
            end: res.end,
            backfill: false,
            roomType: roomType,
        }
        if(roomType == 'dm') {
            props.room_id = room.room_id
        }
        if(roomType == 'thread') {
            props.room_id = room.channel_id
        }
        store.updateEvents(props)

    }).then(() => {
        if(roomType != "dm" && roomType != "thread") {

            // pull events for topics

            let altRoomType;
            if(roomType == 'chat') {
                altRoomType = 'topics'
            } else if(roomType == 'topics') {
                altRoomType = 'chat'
            }

            let popts = {
                room_id: roomID,
                roomType: altRoomType,
            }
            /*
            store.fetchRoomEvents(popts).then(res => {
                console.log(res)
                let props = {
                    room_id: roomID,
                    events: res.chunk,
                    start: res.start,
                    end: res.end,
                    backfill: false,
                    roomType: altRoomType,
                }
                store.updateEvents(props)
            })
            */
        }
    }).then(() => {
        if(thread) {
            getThreadRoom()
        } else {
            fetched = true
            loading = false
            visible = true
        }
    })



})

function getThreadRoom() {

    let roomID = room.thread_in_room_id
    if(!roomID) {
        return
    }
    let opts = {
        room_id: roomID,
    }
    store.fetchRoomEvents(opts).then(res => {
        console.log(res)
        let props = {
            room_id: roomID,
            events: res.chunk,
            start: res.start,
            end: res.end,
            backfill: false,
            isTopics: topics,
        }
        store.updateEvents(props)
    }).then(() => {
        fetched = true
        loading = false
    })
}

$: chatroom = room?.room_type == 'chat'
$: thread = room?.room_type == 'thread'
$: topics = room?.room_type == 'topics'
$: dm = room?.room_type == 'dm'

function activate() {
    if(chatroom || thread || dm) {
        if(!mentionedEventExists) {
            chat.scrollTop = chat.scrollHeight;
        }
    }
    active = true
}

function scrollDown() {
    if(chatroom || thread || dm) {
        setTimeout(() => {
            chat.scrollTop = chat.scrollHeight;
        }, 1)
    }
}

function forceScroll() {
    chat.scrollTop = chat.scrollHeight;
    setTimeout(() => {
        chat.scrollTop = chat.scrollHeight;
    }, 10)
    setTimeout(() => {
        chat.scrollTop = chat.scrollHeight;
    }, 100)
    setTimeout(() => {
        chat.scrollTop = chat.scrollHeight;
    }, 300)
}

function scrollToEvent(e) {
    let el = document.getElementById(e.detail)
    el.scrollIntoView({block: "center"})
}


let fetched = false;

async function getMembers() {
    let endpoint = `${homeServer}/_matrix/client/r0/rooms/${roomID}/joined_members`
    let resp = await fetch(endpoint, {
        headers: {
            'Authorization': `Bearer ${account.matrix_access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

$: servers = $store.accounts.filter(x => x.user_id == $store.active_account)[0].servers


$: rit = `room-${roomAlias}`
$: chit = `ch-${roomAlias}`
$: tps = `tps-${roomAlias}`




let view;

let sh;

let noMore = false;

function loadMore(e) {
    sh = chat?.scrollHeight
    let opts = {
        room_id: e.detail.room_id,
        start: e.detail.from,
        roomType: room.room_type,
    }

    store.fetchRoomEvents(opts).then(res => {
        console.log(res)
        if(res?.chunk?.length < 50) {
            noMore = true
        }
        let props = {
            room_id: roomID,
            events: res.chunk,
            start: res.start,
            end: res.end,
            backfill: true,
            roomType: room.room_type,
        }
        store.updateEvents(props)
    }).then(() => {
        chat.scrollTop =  chat.scrollHeight - sh
    }).then(() => {
    })
}

function loadForward(e) {
    let opts = {
        room_id: e.detail.room_id,
        start: e.detail.from,
        roomType: room.room_type,
        dir: 'f',
    }

    store.fetchRoomEvents(opts).then(res => {
        console.log(res)
        if(res?.chunk?.length < 50) {
            noMore = true
        }
        let props = {
            room_id: roomID,
            events: res.chunk,
            start: res.start,
            end: res.end,
            forward: true,
            roomType: room.room_type,
        }
        store.updateEvents(props)
    }).then(() => {
    })
}

$: events = $store.events[roomID]?.events
$: msgType = `m.room.message`
$: messages = events?.filter(x => x.type == msgType &&
    !x.content?.[`m.new_content`])


$: threadActive = room?.thread?.active

export let replyEvent;

$: noDis = !indicate && !threadView && !isMessageEvent

let newPost;

function newTopic() {
    newPost.activate()
}

</script>


<div id={rit} class="room-root relative" 
    class:hide={!visible} 
    class:show={visible} 
    class:op={visible} 
    class:dis={indicate} 
    class:no-dis={noDis}>

    <div class="room-foundation"
        class:pr1={threadActive}>

        <div id={chit} class="chat-view w-100 scrl" 
            class:no-dis={topics}
            class:is-ch={chatroom || thread || dm} 
            bind:this={chat}>
            {#if fetched}
                <div class=""></div>
                <Events 
                room={room}
                topicsView={false}
                inactive={noDis || topics}
                page={page}
                visible={visible}
                on:scrollDown={scrollDown}
                on:forceScroll={forceScroll}
                on:scrollToEvent={scrollToEvent}
                on:ready={activate}
                noMore={noMore}
                on:replying
                replyEvent={replyEvent}
                on:load-forward={loadForward}
                on:load-more={loadMore}/>
            {/if}
        </div>

        {#if !dm && !thread}
        <div id={tps} class="topics-view w-100 scrl pt4" 
            class:no-dis={chatroom}>
            {#if fetched}
                <div class="gr-default">
                    <div class="create-topic flex mr3 mb3">
                        <div class="gr-center flex-one">
                        </div>
                        <div class="gr-center mr2">
                            <button class="" on:click={newTopic}>Create Topic</button>
                        </div>
                    </div>
                </div>


                <Events 
                room={room}
                inactive={noDis || chatroom}
                topicsView={true}
                page={page}
                visible={visible}
                on:scrollDown={scrollDown}
                on:forceScroll={forceScroll}
                on:scrollToEvent={scrollToEvent}
                on:ready={activate}
                noMore={noMore}
                on:replying
                replyEvent={replyEvent}
                on:load-forward={loadForward}
                on:load-more={loadMore}/>
            {/if}
        </div>
            {#if topics}
                <div class="topic-editor">
                    <NewPost bind:this={newPost} room={room} />
                </div>
            {/if}
        {/if}

    </div>


</div>
    {#if loading}
        <LoadingEvents topics={topics}/>
    {/if}

<style>

.room {
    display: grid;
    width: 100%;
    height: 100%;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
    overflow: hidden;
}

.room-root {
    display: grid;
    width: 100%;
    height: 100%;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
    overflow: hidden;
}

.room-foundation {
    overflow: hidden;
    display: grid;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
}

.topics {
    max-width: 960px;
    width: 100%;
    justify-self: center;
}

.create-topic {
    max-width: 960px;
    width: 100%;
    justify-self: center;
}

@media screen and (max-width: 960px) {
    .topics {
        max-width: 100%;
    }
    .create-topic {
        max-width: 100%;
    }
}


.chat-view {
    overflow: hidden auto;
    scroll-behavior:  auto;
}

.topics-view {
    overflow: hidden auto;
    scroll-behavior:  auto;
    position: relative;
}

.topic-editor {
    position: absolute;
    bottom: 0;
    right: 1rem;
    left: 1rem;
}

.is-ch {
    grid-template-rows: 1fr auto;
    display: grid;
}

.dis {
     display: grid;
     height: 100%;
}
.no-dis {
     display: none;
}
.hide {
    visibility: hidden;
}
.show {
    visibility: visible;
    opacity: 0;
    transition: 0.2s;
}
.op {
    opacity: 1;
}
.loading {
     position: absolute;
     top: 0;
     bottom: 0;
     left: 0;
     right: 0;
     width: 100%;
     height: 100%;
     display: grid;
}
</style>
