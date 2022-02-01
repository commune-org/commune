<script>
import {onMount} from 'svelte'
import { navigate } from 'svelte-navigator'
import { debounce } from '../../../../utils/utils.js'
import { tick } from 'svelte'
import {store} from '../../../../store/store.js'
import {user as userIcon, addSmall, close, check} from '../../../../utils/icons.js'
import isEqual from 'lodash.isequal';

import DirectMessage from './dm/direct-message-item.svelte'
import DMRequest from './dm/dm-request-item.svelte'

import Popup from '../../../../components/popup/popup.svelte'

export let sender;

$: account = $store.accounts?.filter(a => a.user_id == $store?.active_account)?.[0]

$: dms = account?.direct_messages?.filter(x => !x.child)
$: directMessages = dms?.sort((a, b) => (a.origin_server_ts > b.origin_server_ts) ? 1 : -1)

let query = '';
let searchInput;

async function focusInput() {
    await tick()
    searchInput.focus()
}

async function fetchUsers() {
    let data = {
        search_term: query,
        limit: 100,
    }
    let account = $store?.accounts?.filter(x => x.user_id == $store.active_account)[0]
    let endpoint = `${homeServer}/_matrix/client/r0/user_directory/search`
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

async function remoteFetchUsers() {
    let endpoint = `${homeServer}/_matrix/client/r0/profile/${query}`
    let resp = await fetch(endpoint)
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let results = [];

let none = false;

let searching = false;
function searchForUsers() {
    none = false
    showWarning = false
    if(query?.length == 0) {
        return
    }
    searching = true
    /*
    */
    debounce(() => {
        if(query.substring(0,1) == '@') {
            remoteFetchUsers().then(res => {
                if(res) {
                    console.log(res)
                    results = [...results, {
                        user_id: query,
                        display_name: res?.displayname,
                        avatar_url: res?.avatar_url,
                    }]
                }
            }).then(() => {
                searching = false
            })
            return
        }
        fetchUsers().then(res => {
            console.log(res)
            if(res) {
                let ind = res?.results?.findIndex(x => x.user_id == $store.active_account)
                if(ind != -1) {
                    res?.results.splice(ind, 1)
                }
                let indd = res?.results?.findIndex(x => x.user_id.includes(`@commune:`))
                if(indd != -1) {
                    res?.results.splice(ind, 1)
                }
                if(res?.results?.length > 0) {
                    results = res?.results
                    results = results
                } else if(res?.results?.length == 0) {
                    none = true
                }
            }
        }).then(() => {
            searching = false
        })
    }, 500)
}


let selectedUsers = []

$: noneSelected = selectedUsers?.length == 0


function isSelected(user) {
    return selectedUsers.filter(x => x.user_id == user.user_id)[0]
}

function addUser(user) {
    let ind = selectedUsers.findIndex(x => x.user_id == user.user_id)
    if(ind == -1) {
        selectedUsers = [...selectedUsers, user]
        query = ''
        searchInput.value = null
        focusInput()
    } else {
        removeUser(user)
    }
}

function removeUser(user) {
    console.log("deleting", user.display_name)
    let ind = selectedUsers.findIndex(x => x.user_id == user.user_id)
    if(ind != -1) {
        selectedUsers.splice(ind, 1)
        selectedUsers = selectedUsers
        focusInput()
    }
}

function newDMKilled() {
    active = false
    query = ''
    if(searchInput) {
        searchInput.value = null
    }
    results = []
    results = results
    searching = false
    none = false
    selectedUsers = []
    selectedUsers = selectedUsers
}

function reset() {
    selectedUsers = []
    selectedUsers = selectedUsers
}

let container;
function killNewDM() {
    container.kill()
    newDMKilled()
}

let showWarning;
function newDM() {
    if(noneSelected) {
        showWarning = true
        focusInput()
        return
    }
    showWarning = false

    let dms = account.direct_messages?.filter(x => !x.child)

    let sel = []
    selectedUsers.forEach(user => {
        sel.push(user.user_id)
    })
    sel.push($store.active_account)
    sel.sort()

    let exists = false;

    dms.every(dm => {
        let mem = []
        for (const [user_id, _] of Object.entries(dm.members)) {
            mem.push(user_id)
        }
        mem.sort()
        exists = isEqual(mem, sel)
        if(exists) {
            killNewDM()
            navigate(`/messages/${dm.alias}`, {replace:true})
            store.updateActiveHomePage(`/messages/${dm.alias}`)
            store.updateActiveDirectMessages(`/messages/${dm.alias}`)
            return false
        }
        return true
    })
    if(exists) {
        return
    }

    createDM().then(res => {
        console.log(res)
        if(res?.room?.room_id) {
            selectedUsers.forEach(user => {
                account.account_data[user.user_id] = [res?.room?.room_id]
            })
            store.addUserAccountData("m.direct", account.account_data).then(res => {
                console.log(res)
            })
            store.addNewDM(res?.room?.room_id, selectedUsers)
            killNewDM()
        }
    }).then(() => {
    })
}

function buildAvatar(url) {
    console.log(url)
    return `${homeServer}/_matrix/media/r0/download/${url?.substring(6)}`

}

function backspace(e) {
    if(searchInput.value.length > 0) {return}
    if(e.key == 'Backspace' || e.code == 'Backspace') {
        if(selectedUsers?.length > 0) {
            selectedUsers.splice(-1)
            selectedUsers = selectedUsers
            results = []
            results = results
        }
    }
}
function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}

function username(user) {
    return strip(user.user_id)
}


async function createDM() {
    let endpoint = `/server/create`
    let data = {
        dm: true,
        dm_users: []
    };
    selectedUsers?.forEach(user => {
        data.dm_users.push(user.user_id)
    })
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': account.access_token,
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function print(e) {
    e.preventDefault()
    console.log(account?.direct_messages)
}

let active = false;
function activate() {
    active = true
}

onMount(() => {
    window.addNewDM =() => {
        active = true
    }
})

$: dmRequests = account?.dm_requests

</script>

<div class="items-container scrl">

<div class="items flex flex-column">
    {#if dmRequests?.length > 0}
    <div class="ml1 flex ph2 mt1">
        <div class="gr-center lb flex-one" on:contextmenu={print}>
            Requests
        </div>
    </div>

    {#each dmRequests as request (request.room_id)}
        <DMRequest request={request} />
    {/each}

    {/if}

    <div class="ml1 flex ph2 mt2"
        class:mt3={dmRequests?.length > 0}>
        <div class="gr-center lb flex-one" on:contextmenu={print}>
            direct messages
        </div>
        <Popup 
        bind:this={container}
        trigger={"click"}
        on:ready={focusInput}
        on:paste={reset}
        on:killed={newDMKilled}
        shadow={`0 0px 10px rgba(0,0,0,.1)`}
        initActive={active}
        placement={"bottom-start"}
        offset={[0, 10]}>
            <div class="add-icon gr-center mr1" 
                slot="reference">
                {@html addSmall}
            </div>
            <div class="dm-container flex flex-column" slot="content">
                <div class="flex flex-column pa3">
                    <div class="title bold">
                        Start a Direct Message
                    </div>
                    <div class="pt1">
                        You can add multiple users
                    </div>
                </div>
                <div class="input-c flex flex-wrap mh3 relative"
                    on:click|self={focusInput}>
                    {#if selectedUsers?.length > 0}
                        {#each selectedUsers as user, i (user.user_id)}
                            <div class="sel-u pa1 flex ml1  mv1" 
                                on:click|stopPropagation={removeUser(user)}>
                                <div class="flex-one gr-center ph1">
                                    {user.display_name}
                                </div>
                                <div class="discard gr-center">
                                    {@html close}
                                </div>
                            </div>
                        {/each}
                    {/if}
                    <div class="flex flex-one relative">
                        <input
                        bind:this={searchInput}
                        bind:value={query}
                        on:input={searchForUsers}
                        on:keydown={backspace}
                        placeholder={"Search by Username"}/>
                        {#if searching}
                            <div class="spinner-s mr2 gr-center">
                            </div>
                        {/if}
                    </div>
                </div>
                {#if results?.length > 0}
                <div class="users-c flex flex-column ma1 pa1 scrl-s relative">
                    {#each results as user, i (user.user_id)}
                        <div class="user-item flex ph2 pa1"
                            on:click={() => addUser(user)}
                            class:ac-u={results?.length == 1}>
                            {#if user?.avatar_url?.length > 0}
                                <div class="profile-avatar ncol bg-img"
                                    style="background-image: url({buildAvatar(user?.avatar_url)});">
                                </div>
                            {:else}
                                <div class="profile-avatar gr-default">
                                    <div class="log gr-center">
                                        {@html userIcon}
                                    </div>
                                </div>
                            {/if}
                            <div class="flex-one gr-center ml2">
                                <span class="">
                                    <strong>{user.display_name}</strong>
                                </span>
                                <span class="ml1 mute">
                                    {username(user)}
                                </span>
                            </div>
                            <div class="flex gr-default">
                                <div class="s-i gr-default gr-center">
                                <div class="no-dis" class:dis={selectedUsers.filter(x => x.user_id==user.user_id)[0]}>
                                        {@html check}
                                </div>
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
                {/if}
                {#if showWarning}
                    <div class="flex flex-column pt3 ph3 mute">
                        Add at least one user to DM.
                    </div>
                {/if}
                {#if none}
                    <div class="flex flex-column pt3 ph3 mute">
                        No users found.
                    </div>
                {/if}
                <div class="flex flex-column pa3">
                    <button class="pv2" on:click={newDM}>Create DM</button>
                </div>
            </div>
        </Popup>
    </div>

    {#if directMessages.length == 0}
        <div class="mt1 mh3 msg mute">
            You don't have any direct messages. 
            <span class="st-dm" on:click={activate}>
                Start one.
            </span>
        </div>
    {/if}

    <div class="di-c">
        {#each directMessages as item (item.origin_server_ts)}
            <DirectMessage item={item} sender={sender} />
        {/each}
    </div>
</div>
</div>


<style>
.items-container {
    display: grid;
    overflow: hidden auto;
}
.di-c {
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

.add-icon {
    fill: var(--text);
    cursor: pointer;
}
.add-icon:hover {
    fill: var(--white);
}
.st-dm {
    color: var(--text-link);
    cursor: pointer;
}
.st-dm:hover {
    text-decoration: underline;
}
.dm-container {
    border-radius: 8px;
    border: 1px solid var(--background-5);
    background-color: var(--background-3);
    width: 440px;
    color: var(--text);
    font-size: 14px;
    line-height: 1.4;
    display: grid;
    overflow: hidden;
}
.title {
    font-size: 1.2rem;
}
.input-c {
    min-height: 38px;
    border: 1px solid var(--background-1);
    color: var(--text);
    border-radius: 2px;
    background-color: var(--background-1);
}
input {
    min-width: 140px;
    height: 34px;
    width: 100%;
    padding: 0.25rem 0.5rem;
    transition: 0.1s;
}

.searching {
    right: 8px;
    top: 8px;
    position: absolute;
    fill: var(--text);
    width: 16px;
    height: 16px;
    animation-duration: 0.6s;
    animation-name: rotate;
    animation-iteration-count: infinite;
}

@keyframes rotate {
    25% {
        transform: rotate(180deg);
    }
    75% {
        transform: rotate(180deg);
    }

    to {
        transform: rotate(360deg);
    }
}
.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--background-2);
    border-radius: 50%;
    transition: 0.1s;
}
.ncol {
    background-color: transparent;
}
.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
.users-c {
    max-height: 180px;
    overflow: hidden auto;
}
.user-item {
    cursor: pointer;
    border-radius: 6px;
}
.user-item:hover {
    background-color: var(--background-11);
}
.ac-u {
    background-color: var(--background-11);
}
.sel-u {
    background-color: var(--background-3);
    border-radius: 2px;
    cursor: pointer;
}

.discard {
    width: 12px;
    height: 12px;
    fill: var(--text);
}
.sel-u:hover .discard {
    fill: var(--white);
}
.s-i {
    width: 22px;
    height: 22px;
    border-radius: 3px;
    border: 1px solid var(--background-4);
}
.check {
    width: 12px;
    height: 12px;
    fill: var(--text);
}
.no-dis {
    display: none;
}
.dis {
    display: block;
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

.spinner-s {
    height: 16px;
    width: 16px;
}
</style>
