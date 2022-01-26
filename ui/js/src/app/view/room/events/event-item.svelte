<script>
import Editor from '../../../components/editor/editor.svelte'
import Reaction from './reaction/reaction.svelte'
import {onMount, createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher();

import {makeid} from '../../../../utils/utils.js'

//import { twemoji } from '../../../utils/twemoji.js'
import {useLocation, navigate} from 'svelte-navigator'
const location = useLocation()

import { store } from '../../../store/store.js'
import { logo as user } from '../../../switcher/home/logo'
import { thread as threadIcon, addReaction } from '../../../utils/icons.js'

import { formatTime ,timeofDay, formatThreadTime, threadTimeAgo } from '../../../utils/time.js'

import { formatBytes } from '../../../utils/utils'

import EventItemTools from './event-tools/event-item-tools.svelte'

import ImageItem from './components/image/image.svelte'
import AudioPlayer from '../../../components/audio/audio.svelte'
import VideoPlayer from '../../../components/video/video.svelte'
import File from '../../../components/file/file.svelte'
import LinkItem from './components/links/link-item.svelte'

export let index;

export let event;
export let room;
export let rooms;
export let thread;

export let preview;
export let isThread;
export let last;

export let inactive;
export let eventEditing;

export let splitView;

$: username =  strip(event?.sender)
$: replyUsername = replyMember?.user_id ? strip(replyMember?.user_id) : ``

function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}

$: avatar = `${homeServer}/_matrix/media/r0/download/${member?.avatar_url?.substring(6)}`

//$: member = $store.members[room.room_id]?.filter(x => x.user_id == event?.sender)[0]

$: member = $store.allMembers[event?.sender]

$: displayNameExists = member?.display_name?.length > 0
$: avatarExists = member?.avatar_url?.length > 0



$: replyMember = $store.allMembers[event?.content?.['m.relates_to']?.['m.in_reply_to']?.user_id]

$: replyDisplayNameExists = replyMember?.display_name?.length > 0
$: replyAvatarExists = replyMember?.avatar_url?.length > 0

$: replyAvatar = `${homeServer}/_matrix/media/r0/download/${replyMember?.avatar_url?.substring(6)}`

$: contentURL = `${homeServer}/_matrix/media/r0/download/${event?.content?.url?.substring(6)}`



$: name = displayNameExists ? member.display_name : username


$: replyName = replyDisplayNameExists ? replyMember.display_name : replyUsername

$: owner = event?.sender == $store.active_account

$: splitViewEvents = room?.room_type == `chat` ?
    $store.events[room.streams?.['topics']]?.events :
    $store.events[room.streams?.['chat']]?.events 


$: events = splitView ? splitViewEvents : $store.events[room.room_id]?.events

$: msgType = `m.room.message`

$: messages = events?.filter(x => x.type == msgType &&
    !x.content?.[`m.new_content`])


$: prevMsgSameSender = messages?.[index - 1]?.sender == event?.sender 
$: prevMsgDifferentType = messages?.[index - 1]?.type != event?.type
$: nextMsgSameSender = messages?.[index + 1]?.sender == event?.sender 

$: prevMsgOlderThanHour = ((messages?.[index-1]?.unsigned?.age -
    messages?.[index]?.unsigned?.age) > 3600000)

$: relations = event?.unsigned?.['m.relations']

$: isReply = event?.content?.['m.relates_to']?.['m.in_reply_to']

$: edited = relations?.['m.replace'] || event?.edited

$: replyBody = isReply?.body

$: redacted = event?.redacted_because || event?.redacted


function timeAgo(age) {
    let now = new Date()
    let y = new Date() - age
    let postedOn = new Date(y);


    if(postedOn.getDate() == now.getDate()) {
        return `Today at`
    }

    let diff = now.getDate() - postedOn.getDate() 
    if(diff == 1 || diff == -30 || diff == -29) {
        return `Yesterday at`
    }

    return ``
}


$: when = timeAgo(event?.unsigned?.age)

$: threadWhen = threadTimeAgo(hasThreads?.thread_event?.origin_server_ts)
$: threadTime = formatThreadTime(hasThreads?.thread_event?.origin_server_ts)

$: time = formatTime(event?.unsigned?.age)
$: postTime = timeofDay(event?.unsigned?.age)

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

$: if(hasThreads && !fetched)  {
    fetched = true

    let opts = {
        room_id: hasThreads?.room_id,
    }
    /*
    store.fetchRoomEvents(opts).then(res => {
        let props = {
            room_id: hasThreads?.room_id,
            events: res.chunk,
            start: res.start,
            end: res.end,
            backfill: false,
        }
        store.updateEvents(props)
    }).then(() => {
        getMembers(hasThreads?.room_id).then(res => {
            store.updateMembers(hasThreads?.room_id, res.joined)
        }).then(() => {
        })
    })
    */
}



$: roomID = room.room_id


$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: matrix = account?.matrix

$: eventType = `m.room.message`

onMount(() => {
    if(last) {
        dispatch('mounted', index)
    }

    //this is where we send the actual event to the matrix roon
    if(event?.delivered == false) {
        dispatch('scroll', true)
        if(!inactive && eventType) {
            matrix.sendEvent(roomID, eventType, event?.content, event?.transaction_id, (err, res) => {
                if(res?.event_id) {
                }
            });
        }
    }
    if(isThread && name) {
        dispatch('thread-user', name)
    }


    if(isMentionedEvent) {
        dispatch('event-mentioned', id)
    }
})

let root;

let editor;
$: placeholder = `Edit this message.`

let editing = false;

function edit() {
    dispatch('editing', event?.event_id)
    if(last) {
        dispatch('forceScroll', true)
    }
}

$: if(eventEditing == event?.event_id) {
    editing = true
} else {
    editing = false
}

function killEdit() {
    dispatch('editing', null)
}



function reply() {
    dispatch('replying', {
        id: id,
        event: event,
    })
}

export let replyEvent;

$: replying = replyEvent?.event_id === event?.event_id

let body;

function rightClick(e) {
    e.preventDefault()
    console.log(event)
}

/*
let emojiOpts = {
    base: '/static/img/emoji/',
    ext: '.svg',
    className: onlyEmoji ? 'inline-emoji-alt' :  'inline-emoji'
}
*/

$: htmlExists = event?.content?.formatted_body ? event?.content?.formatted_body
    : event?.content?.body ? event?.content?.body : ``


//$: formatted = twemoji.parse(htmlExists, emojiOpts)
$: formatted = formatBody(htmlExists)

function formatBody(content) {
    let div = document.createElement('div')
    div.innerHTML = content
    let spans = div.querySelectorAll('span')
    if(spans?.length > 0) {
        spans.forEach(span => {
            if(span.dataset.userid == $store.active_account) {
                span.classList.add('mentioned')
            }
        })
    }
    return div.innerHTML
}


$: id = `event-${event?.origin_server_ts}`

let editedContent;

function syncEdited(e) {
    editedContent = e.detail
}

function saveEdit(e) {


    let { plain_text, html, length } = editedContent

    if(length == 0) {
        return
    }

    killEdit()

    if(plain_text == event.content.body) {
        return
    }

    let content = {
        "body": ` * ${plain_text}`,
        "msgtype": "m.text",
        "m.new_content": {
            "body": plain_text,
            "formatted_body": html,
            "msgtype": "m.text",
        },
        "m.relates_to": {
            "event_id": event?.event_id,
            "rel_type": "m.replace",
        }
    };

    console.log(content)

    let txnID = makeid(32)

    matrix.sendEvent(roomID, event?.type, content, txnID, (err, res) => {
        console.log(res)
        event.edited = true
    });
    editedContent = null

}


$: roomReactions = events?.filter(x => x.type == 'm.reaction')

$: reactions = roomReactions?.filter(x => x.content?.['m.relates_to']?.event_id == event?.event_id)

$: userReactions = reactions?.filter(x => x.sender == $store.active_account)

$: distinctReactions = distinct(reactions)

let initReactions = 0

let lockReactionCount = false;

function distinct(r) {
    let dis = []
    r?.forEach(reaction => {
        let ind = dis.filter(x => x?.key == reaction?.content?.['m.relates_to']?.key)[0]
        if(!ind) {
            dis.push({
                key: reaction?.content?.['m.relates_to']?.key,
                count: 0,
            })
        } else {
            ind.count = ind.count + 1
        }
    })
    if(!lockReactionCount) {
        initReactions = dis?.length
        lockReactionCount = true
    }
    return dis?.sort((a, b) => (a.count < b.count) ? 1 : -1)
    /*
    let dis = []
    r.forEach(reaction => {
        let ind = dis.findIndex(x => x == reaction?.content?.['m.relates_to']?.key)
        if(ind == -1) {
            dis.push(reaction?.content?.['m.relates_to']?.key)
        }
    })
    return dis
    */
}

$: if(initReactions == 0 && distinctReactions?.length == 1) {
    if(last) {
        dispatch('forceScroll', true)
    }
}

$: isImage = event?.content?.msgtype == 'm.image'
$: isGIF = event?.content?.msgtype == 'commune.gif'
$: isAudio = event?.content?.msgtype == 'm.audio'
$: isVideo = event?.content?.msgtype == 'm.video'
$: isFile = event?.content?.msgtype == 'm.file'
$: isFiles = event?.content?.msgtype == 'm.files'
$: isRecording = event?.content?.msgtype == 'm.recording'
$: isSticker = event?.content?.msgtype == 'commune.sticker'

$: imgURL = event?.content?.url

$: image = `${homeServer}/_matrix/media/r0/download/${imgURL?.substring(6)}`


$: h = event?.content?.info?.h
$: w = event?.content?.info?.w

$: wid = splitView ? 200 : 400

$: base = w < wid ? w : wid

$: width = w > base ? base : w

$: height = (base/w) * h

export let highlightedEvent;

$: isHighlighted = highlightedEvent == event?.event_id


let highlighted = false;

function highlight(e) {
    highlighted = e.detail
    if(highlighted) {
        dispatch('highlighted', event?.event_id)
    } else {
        dispatch('highlighted', null)
    }
}


$: hasCode = event?.content?.formatted?.includes(`<pre><code>`)

$: hasThreads = dm ? account?.direct_messages?.filter(x => x.thread_event?.event_id == event?.event_id)?.[0] : rooms?.filter(x => x.thread_event?.event_id == event?.event_id)?.[0]


$: threadEvents = $store.events[hasThreads?.room_id]?.events
$: threadMessageCount = threadEvents?.filter(x => x.type == 'm.room.message')?.length

$: threadMember = $store.allMembers[hasThreads?.sender]

$: threadDisplayNameExists = threadMember?.display_name?.length > 0
$: threadAvatarExists = threadMember?.avatar_url?.length > 0

$: threadAvatar = `${homeServer}/_matrix/media/r0/download/${threadMember?.avatar_url?.substring(6)}`
$: threadUsername = threadMember?.username
$: threadName = threadDisplayNameExists ? threadMember.display_name : threadUsername

function openThread() {
    store.openThread(room, event, hasThreads, dm)
}


$: serverAlias = $location.pathname.split("/")[1]
$: roomAlias = $location.pathname.split("/")[2]

$: isMessageEvent = $location.pathname.split("/")[3] == `message` &&
    $location.pathname.split("/")[4]?.length > 0

$: isMentionedEvent = isMessageEvent &&
    $location.pathname.split("/")[4] == event?.event_id

$: isTopicEvent = $location.pathname.split("/")[3] == `message` &&
    $location.pathname.split("/")[4]?.length > 0

$: chatroom = room?.room_type == 'chat'
$: topics = room?.room_type == 'topics' && !isMessageEvent
$: thread = room?.room_type == 'thread'
$: dm = room?.room_type == 'dm'

$: newThread = event?.content?.new_thread
$: newThreadTitle = event?.content?.thread_title

$: formatSize = formatBytes(event?.content?.info?.size)


$: hasLinks = event?.content?.links?.length > 0

$: links = event?.content?.links

$: onlyEmoji = stripContent(event?.content?.body)

function stripContent(c) {
    let xx = c?.replaceAll(' ', '')
    let emoji = window.emoji.filter( x => x.unicode == xx)[0]
    return emoji != null & emoji != undefined
}

$: dontHighlight = highlightedEvent != null & highlightedEvent != event?.event_id


function goToTopic() {
    if(chatroom) {
        store.toggleRoomStreamTo(room?.server_id, room?.room_id, 'topics')
    }

    if(splitView) {
        let path = `/${serverAlias}/${roomAlias}/topic/${event.event_id}`
        navigate(path)
        return
    }

    let path = `${$location.pathname}/topic/${event.event_id}`
    navigate(path)
}

$: topicReplies =event?.unsigned?.['m.relations']?.['m.annotation']?.chunk?.filter(x =>x.type=='commune.room.topics.post.reply')?.length


let stickerCanvas;


let reacting = false;
let adc;

function reactToEvent() {
    let opts = {
        room: room,
        event: event,
        container: adc,
        userReactions: userReactions,
        inline: true,
    }
    window.toggleEmojiPicker(opts)
}

$: if($store.emoji?.active && $store.emoji?.event == event?.event_id) {
    reacting = true;
} else {
    reacting = false;
}


$: chatSplit = chatroom && splitView
$: topicsSplit = topics && splitView

function goToMessageEvent() {
    if(splitView) {
        store.toggleRoomStreamTo(room?.server_id, room?.room_id, 'chat')
        let path = `/${serverAlias}/${roomAlias}/message/${event.event_id}`
        navigate(path)
    }
}

$: threadActive = room?.thread?.active

$: eventInThread = room?.thread?.events?.filter(x => x.event_id ==
event.event_id)[0]

$: showEventTools = !editing && 
    event?.delivered && 
    !redacted && 
    !preview && 
    !splitView && 
    !threadActive

function toggleEventInThread(e) {
    if(e.target.tagName == `BUTTON`) {
        return
    }
    store.toggleThreadEvent(room.room_id, event)
    window.focusNewThreadInput()
}

let showSelectUpTo = false;
let showSelectInBetween = false;

$: eventsExist = room?.thread?.events?.length > 0
$: singleEventSelected = room?.thread?.events?.length == 1
$: multipleEventsSelected = room?.thread?.events?.length > 1
$: twoEventsSelected = room?.thread?.events?.length  == 2

$: if(eventsExist) {
    let lastEventSelected = room?.thread?.events[room?.thread?.events?.length-1]
    let before = lastEventSelected?.unsigned?.age > event?.unsigned?.age
    let after = lastEventSelected?.unsigned?.age < event?.unsigned?.age
    if((before || after) && !(before && after)) {
        if(!showSelectInBetween) {
            showSelectUpTo = true
        }
    }
}

$: if(multipleEventsSelected) {
    let firstEventSelected = room?.thread?.events[0]
    let lastEventSelected = room?.thread?.events[room?.thread?.events?.length-1]
    let c1 = firstEventSelected?.unsigned?.age > event?.unsigned?.age
    let c2 = lastEventSelected?.unsigned?.age < event?.unsigned?.age
    if(c1 && c2) {
        showSelectInBetween = true
    }
} else {
    showSelectInBetween = false
}

$: if(showSelectInBetween) {
    showSelectUpTo  =false
}

function selectEventsInBetween() {
    let firstEventSelected = room?.thread?.events[0]
    let lastEventSelected = room?.thread?.events[room?.thread?.events?.length-1]
    let ind0 = messages?.findIndex(x => x.event_id == firstEventSelected.event_id)
    let ind1 = messages?.findIndex(x => x.event_id == lastEventSelected.event_id)
    let eventsInBetween = messages.slice(ind0 + 1 ,ind1)
    store.selectThreadEvents(room.room_id, eventsInBetween)
    window.focusNewThreadInput()
}
function selectEventsUpTo() {
    let firstEventSelected = room?.thread?.events[0]
    let lastEventSelected = room?.thread?.events[room?.thread?.events?.length-1]

    let before = firstEventSelected?.unsigned?.age < event?.unsigned?.age
    let after = lastEventSelected?.unsigned?.age > event?.unsigned?.age

    if(before) {
        let ind0 = messages?.findIndex(x => x.event_id == event.event_id)
        let ind1 = messages?.findIndex(x => x.event_id == lastEventSelected.event_id)
        let eventsUpTo = messages.slice(ind0 ,ind1+1)
        store.selectThreadEvents(room.room_id, eventsUpTo)
        window.focusNewThreadInput()
        return
    }
    if(after) {
        let ind0 = messages?.findIndex(x => x.event_id == firstEventSelected.event_id)
        let ind1 = messages?.findIndex(x => x.event_id == event.event_id)
        let eventsUpTo = messages.slice(ind0 + 1 ,ind1+1)
        store.selectThreadEvents(room.room_id, eventsUpTo)
        window.focusNewThreadInput()
        return
    }
}

function mentionUser() {
    store.mentionUser({
        room_id: roomID,
        username: username,
    })
}

$: fontScaling = $store.settings?.appearance?.fontScaling
$: eventSpacing = $store.settings?.appearance?.eventSpacing
$: isCompact = $store.settings?.appearance?.displayMode === "compact"


</script>


{#if (!splitView && (chatroom || thread || dm)) || (topicsSplit)}
{#if event?.type == 'm.room.message'}


    <div id={id} class="room-event flex flex-column relative" 
    style="font-size:{fontScaling}px;"
    bind:this={root}
    on:click={goToMessageEvent}
    class:re-p={!isCompact}
    class:pointer={splitView}
    class:high={highlighted && isHighlighted}
    class:hov={!highlightedEvent}
    class:room-event-h={dontHighlight}
    class:preview={preview}
    class:mt3={isReply}
    class:editing={editing} 
    class:replying={replying} 
    class:mentioned={isMentionedEvent} 
    on:contextmenu={rightClick}
    class:mb1={!nextMsgSameSender && splitView}
    class:mb3={!nextMsgSameSender && !splitView && !isCompact}>


    {#if isReply}
        <div class="flex reply-to">
            <div class="tick">
            </div>
            <div class="repl-c flex-one flex pr4">
                <div class="">
                    {#if replyAvatarExists}
                        <div class="tiny-avatar bg-img"
                            style="background-image: url({replyAvatar});">
                        </div>
                    {:else}
                        <div class="tiny-avatar gr-default">
                            <div class="tiny-log gr-default gr-center">
                                {@html user}
                            </div>
                        </div>
                    {/if}
                </div>
                <div class="r-u r-s flex">
                    <div class="r-x username">@{replyName}</div>
                    <div class="r-b ml1 flex-one">
                        <div class="r-bb">
                            {replyBody}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}


    <div class="flex relative">

        <div class="ava relative flex flex-column" 
            class:ava-sv={splitView}
            class:pl3={!newThread && !splitView}
            class:pl2={!newThread && splitView}>

            {#if !isCompact}
                {#if !prevMsgSameSender || isReply || prevMsgDifferentType || hasThreads || prevMsgOlderThanHour} 
                    {#if newThread}
                        <div class="thread-icon mb1 gr-center h-100">
                            {@html threadIcon}
                        </div>
                    {:else if avatarExists}
                        <div class="profile-avatar ncol bg-img"
                            class:pa-sv={splitView}
                            style="background-image: url({avatar});">
                        </div>
                    {:else}
                        <div class="profile-avatar gr-default"
                            class:pa-sv={splitView}>
                            <div class="log gr-default gr-center"
                                class:log-sv={splitView}>
                                {@html user}
                            </div>
                        </div>
                    {/if}
                {/if}
            {/if}

            {#if prevMsgSameSender && !isReply && !hasThreads && !prevMsgOlderThanHour}
                    <div class="when-sm gr-default">
                        <div class="gr-center">
                            {postTime}
                        </div>
                    </div>
            {:else if isCompact}
                    <div class="when-sm gr-default op-1">
                        <div class="gr-center">
                            {postTime}
                        </div>
                    </div>
            {/if}


            {#if hasThreads}

                <div class="t-sp">
                </div>
            {/if}
        </div>

        {#if isCompact}
            <div class="gr-center-top username mr2" 
                data-username={username}
                data-userid={event?.sender}
                on:click={mentionUser}>
                <strong>{name}</strong>
            </div>
        {/if}


        <div class="flex flex-column flex-one" 
            class:ev-c={!isCompact}>

            {#if !isCompact}

                {#if (!prevMsgSameSender || isReply || prevMsgDifferentType || hasThreads || prevMsgOlderThanHour) && !newThread}
                    <div class="flex mb1">
                        <div class="gr-center username" 
                            data-username={username}
                            data-userid={event?.sender}
                            on:click={mentionUser}>
                            <strong>{name}</strong>
                        </div>
                        <div class="gr-center ml2 when">
                            {#if !splitView}
                                {when} {time}
                            {:else}
                                {time}
                            {/if}
                        </div>
                    </div>
                {/if}
            {/if}


            {#if newThread}

                <div class="content-body">

                <span class="">{name}</span> <span class="mute">started a thread</span> <span class="">{newThreadTitle}</span> 
                </div>

            {:else if !editing}
            <div class="content-body" 
                class:emj={onlyEmoji && !isCompact}
                class:clmp-1={splitView}
                bind:this={body} 
                class:mute={!event?.delivered}>

                    {#if isImage}

                        <ImageItem
                        splitView={splitView}
                        image={event.content.url}
                        width={event.content.info.w}
                        height={event.content.info.h}/>

                    {:else if isGIF}

                        <ImageItem
                        splitView={splitView}
                        isGIF={true}
                        poster={event.content.info.preview}
                        image={event.content.url}
                        width={event.content.info.w}
                        height={event.content.info.h}/>

                    {:else if isAudio}

                        <AudioPlayer 
                        size={event.content.info.size}
                        title={event.content.body}
                        url={contentURL} />

                    {:else if isRecording}

                        <AudioPlayer 
                        size={event.content.info.size}
                        title={event.content.body}
                        url={contentURL} />

                    {:else if isVideo}

                        <VideoPlayer 
                        title={event.content.body} 
                        size={event.content.info?.size} 
                        height={h}
                        width={w}
                        url={contentURL} />

                    {:else if isFile}
                            <File 
                            title={event.content?.body} 
                            type={event.content.info?.mimetype} 
                            size={event.content.info?.size} 
                            url={contentURL} />
                    {:else if isFiles}
                        <div class="files-items flex flex-column">
                        {#each event?.content?.files as file, i (file.url)}
                            <div class="mb2">
                            {#if file.msgtype == 'm.image'}
                                <ImageItem
                                image={file.url}
                                width={file.info.w}
                                height={file.info.h}/>
                            {:else if file.msgtype == 'm.file'}
                                <File 
                                title={file.info.filename} 
                                type={file.info?.mimetype} 
                                size={file.info?.size} 
                                url={file.url} />
                            {/if}
                            </div>
                        {/each}
                        </div>
                    {:else if isSticker}

                        <div class="sticker-item">
                            <img src={event?.content?.url} />
                        </div>
                    {:else}

                        {#if !redacted}
                            {@html formatted}
                            {#if edited}
                                <span class="edited">(edited)</span>
                            {/if}
                        {:else}
                            <span class="mute">Message Deleted</span>
                        {/if}


                    {/if}

                </div>
            {/if}

            {#if hasLinks}
                {#each links as link, i (link.href)}
                    {#if link?.metadata?.title?.length > 0}
                        <LinkItem 
                        first={i == 0} 
                        last={i == links.length - 1} 
                        link={link} />
                    {/if}
                {/each}
            {/if}



            {#if editing}
                <div class="editing-message mt2">
                    <Editor 
                    placeholder={placeholder}
                    room={room}
                    initial={event.content.formatted_body}
                    inline={true}
                    editing={true}
                    on:cancel-edit={killEdit}
                    on:enter={saveEdit}
                    on:sync={syncEdited}
                    bind:this={editor}/>
                </div>
                <div class="ed-m">
                    escape to <span class="link" on:click={killEdit}>cancel</span> • enter to <span class="link">save</span>
                </div>
            {/if}

        {#if distinctReactions && !preview}
            <div class="event-reactions flex flex-wrap" class:pv1={distinctReactions.length > 0}>
                {#each distinctReactions as reaction, i (i)}
                    <Reaction key={reaction.key} event={event} room={room} />
                {/each}

                {#if distinctReactions?.length > 0}
                <div class="r-ico-h gr-default">
                    <div class="ml1 r-ico gr-default gr-center"
                        class:r-ac={reacting}
                        bind:this={adc}
                        on:click={reactToEvent}>
                        {@html addReaction}
                    </div>
                </div>
                {/if}

            </div>
        {/if}

        </div>

        <div class={splitView ? 'ph2' : 'ph3'}>
        </div>

        {#if showEventTools}
            <div class="event-tools mr3 flex" 
            class:npo={dontHighlight}
            class:op-1={highlighted && isHighlighted}>
                <EventItemTools 
                userReactions={userReactions}
                thread={thread || isThread || hasThreads} 
                hasCode={hasCode} 
                room={room} 
                event={event} 
                on:highlight={highlight} 
                owner={owner} 
                on:edit={edit} 
                on:reply={reply}/>
            </div>
        {/if}
    </div>


    {#if hasThreads && !isThread && !thread && !preview}
        <div class="thread-info-c flex relative">
            <div class="tin-spine pl3">
            </div>
            <div class="thread-info pa2 mt2 mb2 flex flex-column" 
            on:click={openThread}>
                <div class="">
                    <span class="ti-n">{hasThreads?.name}</span>
                    {#if threadMessageCount > 1}
                        <span class="ti-a ml1">{threadMessageCount} messages ›</span>
                    {:else}
                        <span class="ti-a ml1">See thread ›</span>
                    {/if}
                </div>
                <div class="mt1 flex">
                    {#if threadAvatarExists}
                        <div class="tiny-avatar bg-img"
                            style="background-image: url({threadAvatar});">
                        </div>
                    {:else}
                        <div class="tiny-avatar gr-default">
                            <div class="tiny-log gr-default gr-center">
                                {@html user}
                            </div>
                        </div>
                    {/if}
                    <div class="gr-center ml1">
                        <span class="">{threadUsername}</span>
                        {#if newThread}
                            <span class="ml2 wh-s">{when} {time}</span>
                        {:else}
                            <span class="ml2 wh-s">{threadWhen} {threadTime}</span>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    {/if}

    {#if threadActive && !preview}
        <div class="thread-mask flex no-select"
            on:click={toggleEventInThread}>
            <div class="flex-one">
            </div>
            {#if showSelectUpTo && !eventInThread}
                <div class="sel-ic gr-center mr3">
                    <button class="small btc" on:click={selectEventsUpTo}>Up to here</button>
                </div>
            {/if}
            {#if showSelectInBetween && !eventInThread}
                <div class="sel-ic gr-center mr3">
                    <button class="small btc" on:click={selectEventsInBetween}>Select in between</button>
                </div>
            {/if}
            <div class="gr-default gr-center mr3"
                on:click={toggleEventInThread}>
                <div class="att gr-center" 
                    class:atts={eventInThread}
                    on:click={toggleEventInThread}>
                </div>
            </div>
        </div>
    {/if}


</div>


{/if}
{/if}

{#if (!splitView && topics) || (chatSplit)}
    <div id={id} class="topic-event-item pointer"
        class:ph3={!splitView}
        class:ph2={splitView}
        on:click={goToTopic}
        class:o-40={!event?.delivered}
        on:contextmenu={rightClick}>

        <div class="room-event-topic"
            class:pa3={!splitView}
            class:mb3={!splitView}
            class:pa2={splitView}
            class:mb2={splitView}
            class:fresh={event?.fresh}>

            <div class="mr3">
                {#if avatarExists}
                    <div class="profile-avatar ncol bg-img"
                        class:pa-sv={splitView}
                        style="background-image: url({avatar});">
                    </div>
                {:else}
                    <div class="profile-avatar gr-default"
                        class:pa-sv={splitView}>
                        <div class="log gr-default gr-center"
                            class:log-sv={splitView}>
                            {@html user}
                        </div>
                    </div>
                {/if}

            </div>


            <div class="fl-co">
                <div class="topic-title clmp-2">
                    {event?.content?.title}
                </div>
                {#if splitView}
                    <div class="when">
                        {when} {time}
                    </div>
                {/if}
                <div class="topic-replies mute mt2">
                    {#if topicReplies > 1}
                        {topicReplies} replies
                    {:else if topicReplies == 1}
                        1 reply
                    {:else}
                        No Replies
                    {/if}
                </div>
            </div>

            <div class="flex">
            </div>

            {#if !splitView}
                <div class="ml2 when">
                    {when} {time}
                </div>
            {/if}


        </div>

    </div>
{/if}

<style>
.topic-item {
}
.room-event {
}

.re-p {
    padding-top: 0.25rem;
    padding-bottom: 0.25rem;
}

.room-event:hover .event-tools {
    opacity: 1;
}

.room-event-h:hover .event-tools {
    opacity: 0;
}

.npo {
    pointer-events: none;
}

.room-event-topic {
    background-color: var(--room-event-hover);
    border: 1px solid transparent;
    border-radius: 7px;
    cursor: pointer;
    transition: 0.1s;
    display: grid;
    grid-template-columns: auto 1fr auto auto;
}

.fresh {
    animation-name: fade;
    animation-duration: 60s;
}

@keyframes fade {
  from {border: 1px solid var(--green);}
  to {border: 1px solid transparent;}
}

.room-event-topic:hover {
    border: 1px solid var(--background-1);
    box-shadow: 0 4px 14px rgba(0,0,0,.1);
}

:root {
    --room-event-hover: #32353a;
    transition: 0.1s;
}

:root.black {
    --room-event-hover: #000;
    transition: 0.1s;
}


:root.light {
    --room-event-hover: #f8f8f8;
    transition: 0.1s;
}


.hov:hover {
    background-color: var(--room-event-hover);
}


.high {
    background-color: var(--room-event-hover);
}

.editing {
    background-color: var(--room-event-hover);
}

.preview:hover {
    background-color: var(--background-3);
}

.ev-c {
    line-height: 1.3rem;
}

.emj {
    font-size: 48px;
    line-height: 48px;
}

.profile-avatar {
    width: 40px;
    height: 40px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
    cursor: pointer;
}
.pa-sv {
    width: 22px;
    height: 22px;
}
.ncol {
    background-color: transparent;
}

.tiny-avatar {
    width: 16px;
    height: 16px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
    cursor: pointer;
}

.profile-avatar:hover {
}

.username {
    color: var(--white);
    cursor: pointer;
}


.bold {
    font-weight: bold;
}

.username:hover {
    text-decoration: underline;
}
.r-s {
    font-size: 0.9rem;
    color: var(--text-light);
    overflow: hidden;
}
.r-u {
}
.r-x {
    margin-left: 3px;
}
.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.log-sv {
    width: 10px;
    height: 10px;
}

.tiny-log {
    fill: var(--text);
    width: 10px;
    height: 10px;
}

.log.logo-fill {
}

.ava {
    min-width: calc(40px + 2rem);
}
.ava-sv {
    min-width: calc(24px + 1rem);
    margin-top: 4px;
}
.when {
    margin-top: 4px;
    font-size: 0.8em;
    color: var(--text-muted);
}
.wh-s {
    font-size: 0.7em;
    color: var(--text-muted);
}
.when-sm {
    font-size: 0.7em;
    color: var(--text-muted);
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    opacity: 0;
}

.room-event:hover .when-sm {
    opacity: 1;
}

.event-tools {
    border: 1px solid var(--background-1);
    box-shadow: 0 10px 20px rgba(0,0,0,.05);
    border-radius: 7px;
    position: absolute;
    right: 0;
    top: -1rem;
    background-color: var(--background-3);
    opacity: 0;
}
.op-1 {
    opacity: 1;
}

.event-tools:hover {
    box-shadow: 0 6px 6px rgba(0,0,0,.1);
}

.mute {
    color: var(--text-muted);
}
.editing-message {
    background-color:var(--background-9);
    border-radius: 7px;
    width: 100%;
    min-height: 48px;
}
.ed-m {
     font-size: 0.8em;
}
.edited {
     font-size: 0.8em;
     color: var(--text-muted);
}
.replying{
    background-color:var(--background-2);
    -webkit-box-shadow:inset 3px 0px 0px 0px var(--primary);
    -moz-box-shadow:inset 3px 0px 0px 0px var(--primary);
    box-shadow:inset 3px 0px 0px 0px var(--primary);
}

.mentioned{
    background-color:var(--background-2);
    -webkit-box-shadow:inset 3px 0px 0px 0px var(--primary);
    -moz-box-shadow:inset 3px 0px 0px 0px var(--primary);
    box-shadow:inset 3px 0px 0px 0px var(--primary);
}

.content-body {
    word-break: break-word;
}

.content-body p {
    margin-block-start: 0;
    margin-block-end: 0;
}


.link {
    color: var(--blue);
    cursor: pointer;
}

.link:hover {
    text-decoration: underline;
}
.reply-to {
    margin-left: calc(1rem + 20px);
}
.tick {
    width: 33px;
    border-top: 2px solid var(--text-muted);
    border-left: 2px solid var(--text-muted);
    border-bottom: none;
    border-right: none;
    border-radius: 7px 0 0 0 ;
    margin-top: 9px;
    margin-bottom: -1px;
    margin-left: -2px;
    opacity: 0.5;
}
.repl-c {
    margin-left: 5px;
    padding-bottom: 4px;
}

.r-b {
    pointer-events: none;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 1;
    overflow-y: hidden;
    overflow-x: hidden;
    margin-right: 34px;
}

.event-reactions {
}

.t-sp {
}

.t-sp:after {
    border-left: 2px solid var(--text-muted);
    height: calc(100% - 44px);
    content: "";
    position: absolute;
    top: 44px;
    left: 36px;
    width: 30px;
}


.tin-spine {
    width: calc(40px + 2rem);
}

.tin-spine:after {
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 7px;
    content: "";
    height: calc(100% - 35px);
    position: absolute;
    top: 0;
    left: 36px;
    width: 30px;
}

.image-item {
    width: 400px;
}

.thread-info {
    background-color:var(--background-2);
    border-radius: 4px;
    font-size: 0.9em;
    cursor: pointer;
}

.ti-n {
    color: var(--white);
}

.ti-a {
    color: var(--blue);
}

.thread-info:hover .ti-a {
    text-decoration: underline;
}

.topic-title {
    font-size: 1.1rem;
}

.topic-replies {
    font-size: 0.9rem;
}

.sticker-item img {
    width: 200px;
}

.sticker-canvas {
    width: 200px;
    height: 200px;
}

.r-ico-h {
}

.r-ico {
    width: 24px;
    height: 24px;
    fill: var(--text);
    cursor: pointer;
    padding: 0.125rem;
    opacity: 0;
    transition: 0.1s;
}

.r-ico:hover {
    fill: var(--white);
}

.r-ac {
    opacity: 1;
}

.event-reactions:hover .r-ico {
    opacity: 1;
}


.thread-mask {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    cursor: pointer;
}

.att {
    width: 18;
    height: 18;
    border-radius: 50%;
    background-color: var(--background-5);
    border: 2px solid var(--background-5);
    cursor: pointer;
    transition: 0.1s;
}
.atts {
    background-color: var(--green);
}
.thread-mask:hover .att {
    background-color: var(--background-4);
}
.thread-mask:hover .atts {
    background-color: var(--green);
}
.sel-ic {
    opacity: 0;
}
.thread-mask:hover .sel-ic {
    opacity: 1;
}
.btc {
    background-color: #309747;
    padding: 0.125rem 0.25rem;
}
</style>
