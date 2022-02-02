<script>
import { store } from '../../../store/store.js'

import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()
import { useLocation, navigate } from 'svelte-navigator'
import { backAlt as back, people, hamburger, topicsBig, at, hashBig, threadBig, right,splitview } from '../../../utils/icons.js'

import {onMount, tick} from 'svelte'
import tippy from 'tippy.js';

import Search from './search/search.svelte'


const location = useLocation()
$: page = $location.pathname
$: serverAlias = page.split("/")[1]
$: roomAlias = page.split("/")[2]
$: topicPost = page.split("/")[3]?.length > 0 && page.split("/")[3] != 'message'



$: isMobile = $store.isMobile
$: mobileViewToggled = $store.mobileViewToggled

$: showMenu = isMobile

export let showMembers;
export let showSplitView;

export let room;
export let rooms;
export let server;
export let servers;

$: topics = room?.room_type == 'topics'
$: dm = room?.room_type == 'dm'

$: icon = dm? at : topics ? topicsBig : chatroom ? hashBig : thread ? threadBig : ``

function toggleShowMembers() {
    dispatch('toggleShowMembers', true)
}

function toggleSplitView() {
    dispatch('toggleSplitView', true)
}



$: roomName = room?.name
$: roomTopic = room?.topic

$: rooms = servers?.filter(s => s.pathname == `/${server}`)[0]?.rooms

$: threadRoom = rooms?.filter(r => r?.room_id == room?.thread_in_room_id)[0]

//$: threadRoomEvent = room?.thread_event

$: chatroom = room?.room_type == 'chat'
$: thread = room?.room_type == 'thread'


$: switchStream = chatroom ? `topics` : `chat`

onMount(() => {
    tip();
})


let tipIcon;

let tooltip;
let tip = () => {
    if(!tipIcon) {
        return
    }
    tooltip = tippy(tipIcon, {
        content: `Switch to ${switchStream} stream`,
        theme: 'tooltip-html',
        placement: 'bottom',
        offset: [0, 13],
        delay: 0,
        duration: 0,
    });
}


function goToThreadRoom() {
    let pathname = `/${server}${threadRoom.pathname}`
    navigate(pathname)
}

$: hasThreads = rooms?.filter(x => x.thread_in_room_id == room?.room_id)

function toggleMobile() {
    $store.mobileViewToggled = !$store.mobileViewToggled
}

async function toggleStream() {
    store.toggleRoomStream(room?.server_id, room?.room_id)
    store.focusEditorInRoomID = room.room_id
    if(tooltip) {
        tooltip.destroy()
        await tick()
        tip()
    }
}

function goBack() {
    let pathname = `/${server}/${roomAlias}`
    navigate(pathname)
}

</script>

<div class="chat-header" class:sm={showMembers || showSplitView}>
    {#if !topicPost}
    <div class="gr-center w-100 flex">
        <div class="ph3 flex">
            {#if showMenu}
                <div class="ham gr-center pr2" on:click={toggleMobile}>
                    {@html hamburger}
                </div>
            {/if}
            {#if thread}
                <div class="icon-mute gr-center pr2">
                    {@html hashBig}
                </div>
                <div class="name-r gr-center" on:click={goToThreadRoom}>
                    {threadRoom?.name}
                </div>
                <div class="icon-mute gr-center ph2">
                    {@html right}
                </div>
            {/if}
            <div class="icon-mute gr-center pr2">
                {@html icon}
            </div>
            <div class="name gr-center">
                {roomName}
            </div>
        </div>
        <div class=" gr-center">
            {#if roomTopic}
                {roomTopic}
            {/if}
        </div>
        <div class="flex-one">
        </div>

        <Search />

        <div class="icon pointer gr-center mr3" 
            class:sm-ico={showMembers}
            on:click={toggleShowMembers}>
            {@html people}
        </div>

        {#if !dm && !thread}
        <div class="m-icon pointer gr-center mr3" 
            class:m-icon-ac={showSplitView}
            on:click={toggleSplitView}>
            {@html splitview}
        </div>

        <div class="pointer gr-center gr-default mr3 no-select" >
            <div class="stream gr-center gr-default" 
                bind:this={tipIcon}
                class:chat={chatroom}
                class:topics={topics}
                on:click={toggleStream}>
                {#if chatroom}
                    <div class="lb gr-center">
                        topics
                    </div>
                {:else if topics}
                    <div class="lb gr-center">
                        chat
                    </div>
                {/if}
            </div>
        </div>
        {/if}
    </div>
    {:else}
    <div class="gr-center w-100 flex">
        <div class="ph3 flex pointer" on:click={goBack}>
            <div class="back gr-center gr-default mr2">
                {@html back}
            </div>
            <div class="gr-center">
                Go back
            </div>
        </div>
    </div>
    {/if}

</div>

<style>
@font-face {
    font-family: Ginto;
    font-weight: 600;
    src: url(/static/fonts/gi600.woff) format("woff");
}
.chat-header {
    display: grid;
    border-bottom: 1px solid var(--background-1);
}

.sm {
    grid-column: 1/3;
    grid-row: 1;
}

@media screen and (max-width: 1150px) {
.sm {
}

}

.sm-ico {
    fill: var(--white);
}

.name {
    font-family: Ginto;
    font-weight: 600;
}

.name-r {
    font-family: Ginto;
    font-weight: 600;
}
.name-r:hover {
    color: var(--white);
    cursor: pointer;
}

.icon svg {
    color: var(--text);
}

.ham {
    fill: var(--text);
}

.stream {
    background-color: var(--background-1);
    border-radius: 500px;
    padding: 0.25rem;
    min-width: 70px;
    border: 1px solid transparent;
    transition: 0.1s;
}

.chat{
    border: 1px solid var(--green);
}

.topics{
    border: 1px solid var(--blue);
}

.iho {
    fill: var(--text-muted);
}
.lb {
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
    color: var(--text-light);
    transition: 0.1s;
}

.back {
    fill: var(--text-muted);
    width: 24px;
    height: 24px;
}

.m-icon {
    fill: var(--text);
    cursor:pointer;
}

.m-icon:hover {
    fill: var(--white);
}

.m-icon-ac {
    fill: var(--green);
}

.m-icon-ac:hover {
    fill: var(--green);
}

</style>
