<script>
import { store } from '../../../store/store.js'
import RoomItem from './room-item.svelte'
import AddRoom from './add-room/add-room.svelte'
import { downSmall, addSmall } from '../../../utils/icons.js'
import {onMount} from 'svelte'
import tippy from 'tippy.js';

export let server;
export let room;

$: roomAlias = room

$: pathname = `/${server}`


$: servers = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.servers

$: serverID = servers?.filter(server => server.pathname == pathname)[0]?.room_id

$: owner = servers?.filter(server => server.pathname == pathname)[0]?.owner

$: rooms = servers?.filter(server => server.pathname == pathname)[0]?.rooms?.filter(x => x.room_type != "thread")

let adding = false;

function addChannel() {
    if(!owner) {
        return
    }
    adding = true
}

function closeAdd() {
    adding = false
}

let roomsVisible = true;

function toggleRooms() {
    roomsVisible = !roomsVisible
}

function updateRoomType(e) {
}

onMount(() => {
    tip();
})

let container;

let tooltip;
let tip = () => {
    if(!container) {
        return
    }
    tooltip = tippy(container, {
        content: 'Add Channel',
        theme: 'tooltip-html',
        placement: 'top',
        offset: [0, 12],
        delay: 0,
        duration: 0,
    });
}


</script>

{#if rooms}
<div class="server-rooms ">
<div class="sp-r-c scrl flex flex-column pt2">
    {#if rooms}
        <div class="flex ml1">
            <div class="roomtype flex flex-one gr-center" 
                on:click={toggleRooms}>
                <div class="gr-center icon" class:rotate={!roomsVisible}>
                    {@html downSmall}
                </div>
                <div class="lb flex-one">
                    channels
                </div>
            </div>
            {#if owner}
            <div class="gr-center icon pa1 mr2" 
                bind:this={container}
                on:click={addChannel}>
                {@html addSmall}
            </div>
            {/if}
        </div>
        {#if roomsVisible}
            {#each rooms as room (room.room_id)}
                <RoomItem room={room} server={server} roomAlias={roomAlias}/>
            {/each}
        {/if}
    {/if}
</div>
</div>
{/if}

{#if adding}
    <AddRoom server={server} serverID={serverID}
    on:update-room-type={updateRoomType}
    on:kill={closeAdd}/>
{/if}

<style>

.server-rooms {
    overflow: hidden;
    display: grid;
}

.sp-r-c {
    height: 100%;
    overflow: hidden auto;
}

.lb {
    font-family: Ginto;
    text-transform: uppercase;
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--text-light);
    padding: 0.5rem 0.5rem 0.5rem 0.125rem;
    transition: 0.1s;
}
.roomtype {
    transition: 0.1s;
    cursor: pointer;
}
.roomtype:hover .lb{
    color: var(--white);
}
.roomtype:hover .icon{
    fill: var(--white);
}
.rotate {
    transform: rotate(270deg);
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    scroll-behavior: smooth;
}

.scrl:hover {
    scrollbar-color: var(--background-1) transparent;
}

.scrl::-webkit-scrollbar {
  width: 4px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: transparent;
}
.scrl:hover::-webkit-scrollbar-thumb {
  background-color: var(--background-1);
}

</style>
