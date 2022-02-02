<script>
import {store} from '../../../store/store.js'
import {fade, fly} from 'svelte/transition'
import {onMount, createEventDispatcher} from 'svelte'
import Editor from '../../../components/editor/editor.svelte'
import Emoji from './emoji/emoji.svelte'
import GIF from './gif/gif.svelte'
import Stickers from './stickers/stickers.svelte'
import Recording from './recording/recording.svelte'
import Media from './media/media.svelte'
import SendButton from './send-button/send-button.svelte'
import { makeid } from '../../../../utils/utils.js'
import { useLocation } from 'svelte-navigator'

import LinkItem from './links/link-item.svelte'
import FileItem from './files/file-item.svelte'

const location = useLocation()

$: isMessageEvent = $location.pathname.split("/")[3] == `message`

$: indicate = page.pathname == $location.pathname

const dispatch = createEventDispatcher()

export let room;
export let page;
export let threadView;

$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics'
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

$: roomID = room?.room_id
$: sender = $store.accounts.filter(x => x.user_id == $store.active_account)[0].user_id

let editor;
let editorContent = {plain_text: null, html: null, length: 0};

function syncContent(e) {
    editorContent = e.detail
}

$: lengthExceeded = editorContent?.length > 2000
$: exceededBy = editorContent?.length - 2000

$: noInput = editorContent?.length == 0

export let replying;
export let replyEvent;

$: if(replyEvent) {
    console.log("scrollll")
    setTimeout(() => {
        window.syncEventsPosition(room?.room_id)
    },10)
}


let LOCKED = false;

let showLengthWarning = false

function killLengthWarning() {
    showLengthWarning = false
    editor.focusEditor()
}


function send(e) {
    if(LOCKED) {
        return
    }
    if(lengthExceeded) {
        showLengthWarning = true
        return
    }

    let tempid = makeid(16)

    let {plain_text, html, length} = editorContent

    if(length == 0 && files?.length == 0) {
        return
    }

    let content = {
        "body": plain_text || '',
        "msgtype": "m.text",
    };

    if(html?.length > 0) {
        content["formatted_body"] = html
    }

    if(links?.length > 0) {
        content["links"] = links
    }

    if(files?.length == 1) {
        content.msgtype = files[0]?.msgtype
        content.body = files[0].info?.filename
        content.info = files[0].info
        content.url = files[0].url
        console.log("just a single file", content)
    } else if(files?.length > 1) {
        content.files = []
        files.forEach(file => {
            if(file?.url?.length > 0) {
                content.files.push(file)
            }
        })
        content.body = "files"
        content.msgtype = "m.files"
        console.log("multiple files", content)
    }


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

    if(replying && replyEvent) {
        discardReply()
    }
    editorContent = {plain_text: null, html: null, length: 0}
    links = []
    links = links
    files = []
    files = files
    dispatch('scroll-to-bottom', true)
}

function discardReply() {
    dispatch('discardReply', true)
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

function insert(e) {
    editor.insertEmoji(e.detail)
    editor.focusEditor()
}


$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

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

function createThread() {
    editor.reset()
}

$: plt = dm ? `@` : `#`

$: placeholder = `Message ${plt}${room?.name}`


let isRecording = false;

function toggleRecording() {
    isRecording = !isRecording
    if(!isRecording) {
        editor.focusEditor()
    }
}

$: if(replying) {
    editor.focusEditor()
}

let files = [];
function filesReady(e) {
    files = e.detail
    console.log(files)
    editor.focusEditor()
}

function discardFile(e) {
    const id = e.detail
    let ind = files.findIndex(x => x.id == id)
    if(ind != -1 ) {
        files.splice(ind, 1)
        files = files
        editor.focusEditor()
    }
}

function fileUploaded(e) {
    const {id, url} = e.detail
    let ind = files.findIndex(x => x.id == id)
    if(ind != -1) {
        files[ind].url = url
    }
    LOCKED = false
    console.log("FILES UPDATED", files)
}

let links = []

function newLink(e) {
    let link = e.detail
    let ind = links.findIndex(x => x.href == link.href)
    if(ind == -1) {
        links.push(link)
        links = links
    }
    console.log("LINKS ARE", links)
}

function discardLink(e) {
    const href = e.detail
    let ind = links.findIndex(x => x.href = href)
    if(ind != -1 ) {
        links.splice(ind, 1)
        links = links
    }
}

function updateMetadata(e) {
    const {href, metadata} = e.detail
    let ind = links.findIndex(x => x.href == href)
    if(ind != -1 ) {
        links[ind].metadata = metadata
    }
    LOCKED = false
}

function fetchingMetadata(e) {
    LOCKED = true
}


function keyPress(e) {
    if(e.key == 'Escape' || e.code == 'Escape') {
        killLengthWarning()
    }
}
$: if(showLengthWarning) {
    document.addEventListener('keydown', keyPress)
} else {
    document.removeEventListener('keydown', keyPress)
}

$:id = `inv-${room.alias}`

$: showMessageButton = $store.settings.accessibility.showMessageButton

</script>

<div id={id} class="input-view" 
    class:recording={isRecording}
    class:smb={showMessageButton}
    class:no-dis={!indicate && !threadView && !isMessageEvent}>

    {#if lengthExceeded}
        <div class="exceeded pa2">
            -{exceededBy}
        </div>
    {/if}

    {#if !isRecording}
    <Media 
    on:files-ready={filesReady}
    editorContent={editorContent} 
    room={room} 
    page={page}
    on:createThread={createThread}/>
    {/if}

    <div class="editor-area">
        <Editor 
        recording={isRecording}
        on:sync={syncContent}
        page={page}
        on:enter={send}
        threadView={threadView}
        placeholder={placeholder}
        focus={indicate} 
        room={room}
        on:new-link={newLink}
        bind:this={editor}/>

        {#if links?.length > 0}
            <div class="">
                {#each links as link (link.href)}
                    <LinkItem 
                    href={link.href} 
                    room={room} 
                    on:metadata-fetched={updateMetadata}
                    on:fetching-metadata={fetchingMetadata}
                    on:discard={discardLink}/>
                {/each}
            </div>
        {/if}
        {#if files?.length > 0}
            <div class="">
                {#each files as file, i (file.id)}
                    <FileItem 
                    file={file} 
                    room={room} 
                    last={i == files.length - 1}
                    on:file-uploaded={fileUploaded}
                    on:discard={discardFile}/>
                {/each}
            </div>
        {/if}
    </div>

    <Recording 
    on:toggleRecording={toggleRecording}
    recording={isRecording}
    room={room} 
    page={page}/>

    {#if !isRecording}
        <GIF on:selected={insertGIF}/>
        <Stickers room={room} on:selected={insertGIF}/>
        <Emoji on:insert={insert}/>

        {#if showMessageButton}

            <SendButton noInput={noInput} on:send={send}/>

        {/if}

    {/if}

</div>

{#if showLengthWarning}
    <div class="mask gr-default" 
        on:click={killLengthWarning}
        transition:fade="{{duration: 100}}">
        <div class="modal gr-center" 
            in:fly="{{ y: -200, duration: 100 }}">

            <div class="too-large pa2">
                <div class="flex flex-column tl-i pa3">
                    <div class="gr-center bold">
                        Your message is too long.
                    </div>
                    <div class="gr-center pt2">
                        Max length is 2000 characters.
                    </div>
                </div>
            </div>

        </div>
    </div>

{/if}


<style>
.input-view {
    display: grid;
    grid-template-columns: auto 1fr auto auto auto auto;
    height: 100%;
}

.smb {
    grid-template-columns: auto 1fr auto auto auto auto auto;
}

.recording {
}

.no-dis {
    display: none;
}
.editor-area {
    display: grid;
    grid-template-rows: auto auto;
}
.exceeded {
    background-color: var(--background-3);
    border-radius: 8px;
    position: fixed;
    bottom:2rem;
    right: 2rem;
    color: red;
}
.mask {
    transition: 0.3s;
    position: fixed;
    z-index: 1000;
    width: 100%;
    height: 100%;
    background-color: var(--mask);
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
}

.modal {
    transition: 0.2s;
    box-shadow: 0 30px 60px rgba(0,0,0,.1);
}

.too-large {
    background-color: hsl(359,calc(var(--saturation-factor, 1)*66.7%),54.1%);
    border-radius: 12px;
    color: var(--white);
    width: 310px;
}

.tl-i {
    border-radius: 10px;
    border: 2px dashed var(--text);
}
</style>
