<script>
import { store } from '../../store/store.js'
import {onMount, tick} from 'svelte'
import { accountSettings, back } from '../../utils/icons.js'
import { logo as userIcon } from '../../switcher/home/logo'

import AccountItem from './account-item.svelte'

import { fade, fly } from 'svelte/transition'

import Popup from '../../components/popup/popup.svelte'

$: avatarExists = user?.avatar_url?.length > 0

$: avatar = `${homeServer}/_matrix/media/r0/download/${user?.avatar_url?.substring(6)}`

$: user = $store.accounts.filter(x => x.user_id == $store.active_account)[0]

$: accounts = $store.accounts

$: matrix = user?.matrix

$: username = user?.username

$: status = user?.status

let popup;
onMount(() => {
    //initMenu()
})

async function fetchStatus(presence, status_msg) {
    let endpoint = `${homeServer}/_matrix/client/r0/presence/${user?.user_id}/status`
    let account = $store?.accounts?.filter(account => account.user_id == $store.active_account)[0]
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

function setStatus(status) {
    popup.kill()
    let data = {
        presence: status,
        status_msg: "",
    }
    matrix.setPresence(data, (err, res) => {
        console.log(res)
    })
}

let active = false;

function toggle() {
    active = !active
    if(!active) {
        adding = false
    }
}

$: if(active) {
    document.addEventListener('keydown', escape)
} else {
    document.removeEventListener('keydown', escape)
}
function escape(e) {
    if(e.key == 'Escape') {
        if(adding) {
            cancelAdding()
            return
        }
        active = false
        reset()
    }
}


let adding = false;

function addAccount() {
    adding = true
    focusUsernameInput()
}

function cancelAdding() {
    resetAdding()
}

function resetAdding() {
    adding = false
    newUsername = null
    newPassword = null
    loginError = false
}

let newUsername;
let usernameInput;
let newPassword;
let passwordInput;

async function focusUsernameInput() {
    await tick()
    usernameInput.focus()
}
async function focusPassswordInput() {
    await tick()
    passwordInput.focus()
}

async function fetchCredentials() {
    let endpoint = `/api/v0/login`
    let data = {
        username: newUsername,
        password: newPassword,
    };
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        credentials: "same-origin",
        headers:{
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let validating = false;
let loginError = false;

function addNewAccount() {
    if(usernameInput.value.length == 0 ) {
        usernameInput.focus()
        return
    }

    let exists = $store.accounts.filter(x => x.username == newUsername)[0]
    if(exists) {
        alert(`You're already logged into that account.`)
        return
    }
    if(passwordInput.value.length == 0 ) {
        passwordInput.focus()
        return
    }
    validating = true
    fetchCredentials().then(res => {
        validating = false
        console.log(res)
        if(res?.error) {
            loginError = true
        }
        if(res?.authenticated) {
            store.addNewAccount({
                identity: res?.identity,
            })
            resetAdding()
        }
    })
}

function onUsernameKeyPress(e) {
    if(e.key == 'Enter') {
        if(passwordInput.value.length == 0 ) {
            passwordInput.focus()
            return
        }
        addNewAccount()
    }
}


function onPasswordKeyPress(e) {
    if(e.key == 'Enter') {
        addNewAccount()
    }
}

</script>


<div class="profile-container fl-co">

<div class="sidebar-profile flex pr1">


    <Popup
    bind:this={popup}
    trigger={"click"}
    shadow={`0 0px 10px rgba(0,0,0,.1)`}
    borderRadius={`4px`}
    placement={"top-start"}
    offset={[0, -14]}>
        <div class="profile-avatar pointer gr-center gr-default relative"
        slot="reference"
        class:darj-bg={!avatarExists}>
            {#if !avatarExists}
                <div class="logo gr-center">
                    {@html userIcon}
                </div>
            {:else}
                <div class="avatar bg-img"
                style="background-image: url({avatar});">
                </div>
            {/if}
            <div class="presence">
            </div>
        </div>
        <div class="menu-container pa2" slot="content">
            <div class="menu fl-co no-select">
                <div class="menu-item flex" 
                    on:click={() => setStatus('online')}>
                    <div class="gr-center pr2">
                        <div class="online"></div>
                    </div>
                    <div class="">
                        Online
                    </div>
                </div>
                <div class="sep"></div>
                <div class="menu-item flex" 
                    on:click={() => setStatus('offline')}>
                    <div class="gr-center pr2">
                        <div class="idle"></div>
                    </div>
                    <div class="">
                        Offline
                    </div>
                </div>
                <div class="menu-item flex" 
                    on:click={() => setStatus('unavailable')}>
                    <div class="gr-center pr2">
                        <div class="offline"></div>
                    </div>
                    <div class="">
                        Unavailable
                    </div>
                </div>
                <div class="sep"></div>
                <div class="menu-item flex">
                    <div class="gr-center flex-one">
                        Set a custom status
                    </div>
                </div>
            </div>
        </div>
    </Popup>


    <div class="gr-center flex flex-column flex-one">
        <div class="usn gr-center-start pl2">
            {username}
        </div>
    </div>

    <div class="ac-i icon pa1 pointer gr-default" 
        class:flip-v={active}
        on:click={toggle}>
        {@html accountSettings}
    </div>
</div>


{#if active}
<div class="mask gr-default" 
    on:click|self={toggle}
    transition:fade="{{duration: 100}}">
    <div class="modal gr-center"
        in:fly="{{ y: -200, duration: 100 }}">

        {#if !adding}
            <div class="accounts pa3 fl-co">

                <div class="lbl mb3">
                    accounts
                </div>
                <div class="account-items fl-co">
                    {#each accounts as account, i (account.user_id)}
                        <AccountItem account={account} />
                    {/each}
                </div>

                <div class="pt3">
                    <button class="" on:click={addAccount}>
                        Add another account
                    </button>
                </div>

            </div>
        {/if}

        {#if adding}
            <div class="accounts pa3 fl-co">

                <div class="lbl mb3">
                    Add new account
                </div>
                <div class="account-items fl-co mt3">

                    <div class="lbl mb2">
                        username
                    </div>
                    <div class="input-holder">
                        <input type="text"
                        bind:this={usernameInput}
                        bind:value={newUsername}
                        on:keypress={onUsernameKeyPress}
                        placeholder="username"
                        />
                    </div>

                    <div class="lbl mb2 mt3">
                        password
                    </div>
                    <div class="input-holder">
                        <input type="password"
                        bind:this={passwordInput}
                        bind:value={newPassword}
                        on:keypress={onPasswordKeyPress}
                        placeholder="password"
                        />
                    </div>

                    <div class="mt3 warn" style="min-height: 20px;">
                        {#if loginError}
                            Username or password did not match.
                        {/if}
                    </div>

                </div>

                <div class="pt3 flex">
                    <div class="gr-center pointer" on:click={cancelAdding}>
                        {@html back}
                    </div>
                    <div class="gr-center flex-one">
                    </div>
                    {#if validating}
                        <div class="spinner-s mr2">
                        </div>
                    {/if}
                    <div class="gr-center">
                        <button class="" 
                            disabled={validating}
                            on:click={addNewAccount}>
                            add account
                        </button>
                    </div>
                </div>

            </div>
        {/if}
    </div>
</div>
{/if}



</div>


<style>
@font-face {
    font-family: Whitney;
    font-weight: 600;
    src: url('/static/fonts/wh600.woff') format('woff');
}

.usn {
    font-weight: 600;
}
:root {
    --online: #3ba55d;
}
.profile-container {
    background-color: var(--background-5);
}
.sidebar-profile {
    padding: 0.5rem;
    background-color: var(--background-5);
}
.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
}
.profile-avatar:hover {
    opacity: 0.8;
}
.dark-bg {
    background-color: var(--dark);
}
.presence {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    border: 4px solid var(--background-5);
    background-color: var(--online);
    position: absolute;
    bottom: -2px;
    right: -2px;
}

.online {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--online);
}

.idle {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: red;
}

.offline {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--background-3);
}


.logo {
    width: 16px;
    fill: var(--white);
}
.avatar {
    width: 100%;
    height: 100%;
    border-radius: 50%;
}
.menu-container {
    width: 222px;
    background-color: var(--menu);
}
.menu-item{
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
    font-size: 0.9rem;
}
.menu-item:hover{
    color: white;
    background-color: var(--background-2);
}

.ico {
    fill: var(--text-muted);
    cursor: pointer;
    height: 30px;
    width: 30px;
    border-radius: 4px;
}
.ico:hover {
    fill: var(--white);
    background-color: var(--background-3);
}
.username {
    color: var(--text);
    font-weight: bold;
}
.sep {
    border-bottom: 1px solid var(--background-2);
    margin-top: 0.25rem;
    margin-bottom: 0.25rem;
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
    background-color: var(--background-1);
    width: 440px;
    min-height: 300px;
    border-radius: 7px;
    transition: 0.2s;
    box-shadow: 0 30px 60px rgba(0,0,0,.1);
    display: grid;
}
.lbl {
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
    color: var(--text);
}
.accounts {
    display: grid;
    grid-template-rows: auto 1fr auto;
}
button {
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
}

input {
    width: 100%;
    border: 1px solid var(--background-1);
    background-color: var(--background-2);
    border-radius: 5px;
    padding: 0.75rem;
    transition: 0.1s;
}

input:focus {
    border: 1px solid var(--green);
}

.ac-i {
    width: 26px;
}

</style>
