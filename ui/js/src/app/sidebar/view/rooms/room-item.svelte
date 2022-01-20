<script>
import {onMount} from 'svelte'
import { navigate, useLocation } from 'svelte-navigator'
import { store } from '../../../store/store.js'
import { hash, topics, close } from '../../../utils/icons.js'

const location = useLocation()

$: path = $location.pathname

export let server;
export let room;
export let roomAlias;

$: spc = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.servers

$: thisServer = spc?.filter(x => x.room_id == room.server_id)[0]

$: activePath = `${thisServer.pathname}${room.pathname}`

$: indicate = path.substring(0, activePath.length) == activePath

$: pathname = `/${server}${room.pathname}`

$: if(indicate) {
    document.title = room.name
    localStorage.setItem('last-active-page', pathname)
}


$: joined = room?.joined

let threadIndicate = (thread) => {
    let ind = path == `${thisServer.pathname}${thread.pathname}`
    if(ind) {
        document.title = thread.name
    }
    return ind
}

let go =() => {
    navigate(pathname)
    //store.updateActiveRoom(`/${server}`, pathname)
    store.updateActiveRoom(`/${server}`, pathname)
    store.updateActiveRooms(`/${server}`, pathname)
    /*
    if($store.isMobile && $store.mobileViewToggled) {
        $store.mobileViewToggled = false
    }
    */
}

let goToThread =(thread) => {
    if(!thread.joined) {
        store.joinRoom(thread.room_id)
    }
    //let path = `${pathname}${thread.pathname}`
    let path = `/${server}${thread.pathname}`
    navigate(path)
    document.title = thread.name
    store.updateActiveRoom(`/${server}`, path)
    store.updateActiveRoom(`/${server}`, path)
}

onMount(() => {
    if(indicate) {
        store.updateActiveRoom(`/${server}`, pathname)
        store.updateActiveRooms(`/${server}`, pathname)
    }
    if(!joined) {
        //store.joinRoom(room.room_id)
    }
})

function info(e) {
    e.preventDefault()
    console.log(room)
}

function tinfo(e) {
    e.preventDefault()
    console.log(e)
}

$: children = thisServer.rooms?.filter(x => x?.thread_in_room_id == room.channel_id)

$: hasChildren = children?.length > 0

$: single = children?.length == 1

function last(i) {
    return i == children?.length - 1
}


function deleteRoom(roomID) {
    deleteFetch(roomID).then((res) => {
        console.log(res)
    })
}
async function deleteFetch(roomID) {
    let data = {
        room_id: roomID,
    }
    let endpoint = `/server/purge`
    let account = $store?.accounts?.filter(account => account.user_id == $store.active_account)[0]
    console.log(account)
    let resp = await fetch(endpoint, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
            'Authorization': `${account.access_token}`
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

$: chatroom = room?.room_type == 'chat'

$: icon = chatroom ? hash : topics

</script>



<div class="room-item flex flex-column" on:click={go} on:contextmenu={info}>
    <div class="room-item-container flex" class:highlight={indicate}>
        <div class="icon gr-default pr2" class:highlight={indicate}>
                {@html icon}
        </div>
        <div class="clmp-1" >
            {room.name}
        </div>
    </div>
</div>
    {#if hasChildren}
    {/if}

<style>
.room-item {
    padding: 0 0.5rem;
}
:root {
    --ri-ni: #40444b;
}
:root.light {
    --ri-ni: #f2f3f5;
}


.room-item-container {
    padding: 0.5rem 0.75rem;
    border-radius: 5px;
    cursor: pointer;
    margin-bottom: 3px;
}
.room-item-container:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.icon {
}

.room-item-container:hover .icon {
    fill: var(--white);
}

.highlight {
    background-color: var(--background-8);
    color: var(--white);
    fill: var(--white);
}
.highlight:hover {
    background-color: var(--background-8);
}

.spine {
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-left: 2px solid var(--text-muted);
    position: relative;
}

.spine:after {
    background: var(--text-muted);
    content: "";
    height: 2px;
    position: absolute;
    top: 16px;
    width: 10px;
}

.spine-last {
    margin-top: -15px;
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-radius: 0 0 0 4px;
}

.spine-last:after {
    top: 0;
    margin-top: 0;
    margin-left: 0;
    width: 10px;
    height: 18px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
    content: "";
    position: absolute;
}

.spine-single {
    margin-left: 15px;
    width: 0.5rem;
    height: 20px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
}


.thread-name {
    border-radius: 5px;
    cursor: pointer;
    color: var(--text-light);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
}

.thread-name:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.thread-name-ac {
    border-radius: 5px;
    cursor: pointer;
    color: var(--white);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
    background-color: var(--background-6);
}

.thread-name-ac:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.room-item-threads {
    padding: 0 0.5rem;
}
</style>
