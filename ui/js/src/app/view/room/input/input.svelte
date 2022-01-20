<script>
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()
import { useLocation } from 'svelte-navigator'

import { makeid } from '../../../../utils/utils.js'

import { store } from '../../../store/store.js'
import { close } from '../../../utils/icons.js'
import Editor from '../../../components/editor/editor.svelte'
import Emoji from './emoji/emoji.svelte'
import GIF from './gif/gif.svelte'
import Recording from './recording/recording.svelte'
import Media from './media/media.svelte'

import InputView from './input-view.svelte'

const location = useLocation()

export let threadView;
export let page;



$: page = threadView ? page?.pathname : $location.pathname

$: server = page.split("/")[1]
$: roomAlias = page.split("/")[2]

$: isDirectMessage = server == 'messages'

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]
$: servers = account?.servers

$: rooms = servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms

$: dms = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.direct_messages

$: dmRoom = dms?.filter(x => x.alias == roomAlias) [0]

$: room = isDirectMessage ? dmRoom : rooms?.filter(r => r.pathname == `/${roomAlias}`)[0]

$: roomID = room?.room_id


export let replying;
export let replyEvent;


$: indicate = page == window.location.pathname

$: activeDMs = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.active_direct_messages

$: activeRooms = isDirectMessage ? activeDMs : account.active_rooms

$: placeholder = `Message #${room?.name}`

$: sender = $store.accounts.filter(x => x.user_id == $store.active_account)[0].user_id

let editor;

let editorContent;

function syncContent(e) {
    editorContent = e.detail
}

function send(e) {
    let tempid = makeid(16)

    if(!editorContent?.plain_text) {
        return
    }

    let {plain_text, html, length} = editorContent

    let content = {
        "body": plain_text,
        "formatted": html,
        "msgtype": "m.text",
    };


    if(replying && replyEvent) {
        content.format = "org.matrix.custom.html"
        /*
        let repb = replyEvent.content.body
        if(repb?.length >= 100) {
            repb = repb.substring(0,99)
        }
        */
        content['m.relates_to'] = {
            'm.in_reply_to':  {
                event_id: replyEvent?.event_id,
                user_id: replyEvent?.sender,
                body: replyEvent?.content.body,
            }
        }
    }
    console.log(content)


    let ts = new Date()
    store.addEventToRoom(roomID, {
        "type": "m.room.message",
        "room_id": roomID,
        "sender": sender,
        "content": content,
        "origin_server_ts": ts,
        "unsigned": {
            "age": 0,
        },
        "event_id": tempid,
        "user_id": sender,
        "age": 0,
        "delivered": false,
        "transaction_id": tempid,
    })

    /*
    matrix.sendEvent(roomID, "m.room.message", content, "", (err, res) => {
        console.log(res);
    });
    */
    if(replying && replyEvent) {
        discardReply()
    }
    editorContent = null
    dispatch('scroll-to-bottom', true)
}

function insertGIF(e) {

    let tempid = makeid(16)

    let {info, url} = e.detail

    let content = {
        "body": url,
        "info": info,
        "url": url,
        "msgtype": "commune.gif",
    };

    let ts = new Date()
    store.addEventToRoom(roomID, {
        "type": "m.room.message",
        "room_id": roomID,
        "sender": sender,
        "content": content,
        "origin_server_ts": ts,
        "unsigned": {
            "age": 0,
        },
        "event_id": tempid,
        "user_id": sender,
        "age": 0,
        "delivered": false,
        "transaction_id": tempid,
    })

    dispatch('scroll-to-bottom', true)
}

$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

function discardReply() {
    dispatch('discardReply', true)
}

$: if(replyEvent) {
    console.log("YIKES")
    console.log("YIKES")
    console.log(replyEvent)
    console.log(rm,room?.server_id)
    console.log("YIKES")
    console.log("YIKES")
}

$: rm = $store.members

$: replyingTo = replyEvent ?
    Object.entries($store.members[room?.server_id])?.filter(x => x.user_id ==
        replyEvent?.sender)[0] :  null

$: displayName = replyingTo?.display_name

$: name = displayName?.length > 0 ? displayName : strip(replyEvent?.sender)

function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}
function insert(e) {
    editor.insertEmoji(e.detail)
    editor.focusEditor()
}

$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

function createThread() {
    editor.reset()
}

function uploaded(e) {
    console.log(e.detail)
    let tempid = makeid(16)

    let {info, msgtype, room_id, url} = e.detail

    let content = {
        "body": url,
        "info": info,
        "url": url,
        "msgtype": msgtype,
    };

    let ts = new Date()
    store.addEventToRoom(room_id, {
        "type": "m.room.message",
        "room_id": room_id,
        "sender": sender,
        "content": content,
        "origin_server_ts": ts,
        "unsigned": {
            "age": 0,
        },
        "event_id": tempid,
        "user_id": sender,
        "age": 0,
        "delivered": false,
        "transaction_id": tempid,
    })

    dispatch('scroll-to-bottom', true)
}

</script>

<div class="chat-container no-select" 
    class:rep={replying}
    class:no-dis={topics}>

    {#if replying}
        <div class="replying flex w-100">
            <div class="gr-start-center flex-one pl3">
                replying to <strong>{name}</strong>
            </div>

            <div class="close-icon gr-center mr3" on:click={discardReply}>
                {@html close}
            </div>
        </div>
    {/if}

    <div class="chat-input scrl" class:re-ch={replying}>


        <div class="input-container">
            {#if !threadView}
                {#each activeRooms as p, i (p.pathname)}

                        <InputView
                        replying={replying}
                        replyEvent={replyEvent}
                        threadView={threadView}
                        on:discardReply={discardReply}
                        editorContent={editorContent} 
                        page={p}
                        room={room} />

                {/each}

            {:else}

                        <InputView
                        replying={replying}
                        replyEvent={replyEvent}
                        threadView={threadView}
                        on:discardReply={discardReply}
                        editorContent={editorContent} 
                        page={page}
                        room={room} />


            {/if}
        </div>
    </div>
</div>

<style>
.chat-container {
    padding-right: 1rem;
    padding-left: 1rem;
    padding-bottom: 1.25rem;
    display: grid;
    grid-template-columns: 1fr;
    overflow:hidden;
}


:root {
    --chat-input: #40444b;
}
:root.light {
    --chat-input: #f2f3f5;
}


.chat-input {
    min-height: 48px;
    background-color: var(--chat-input);
    width: 100%;
    max-height: 40vh;
    overflow: hidden scroll;
    border-radius: 7px;
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: 1fr auto auto;
}
.re-ch {
    border-radius: 0 0 7px 7px;
}
.holder:hover .icon{
    fill: var(--white);
}
.replying {
    background-color: var(--background-1);
    height: 32px;
    border-radius: 7px 7px 0 0 ;
}
.no-dis {
    display: none;
}
.mt {
    margin-top: 13px;
}
.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--background-1) transparent;
    scroll-behavior: auto;
}

.scrl::-webkit-scrollbar {
  width: 4px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: var(--background-1);
}

.close-icon {
    height: 22px;
    width: 22px;
    cursor: pointer;
    fill: var(--text);
}

.close-icon:hover {
    fill: var(--white);
}
</style>
