<script>
import { onMount } from 'svelte'
import { closeBig } from '../utils/icons.js'
import Nav from './nav/nav.svelte'

import Account from './components/account/account.svelte'
import Appearance from './components/appearance/appearance.svelte'
import Accessibility from './components/accessibility/accessibility.svelte'
import Language from './components/language/language.svelte'

let active = true;

let kill =() => {
    active = false;
}

onMount(() => {
    window.toggleSettings = () => {
        active = !active
    }
})

$: if(active) {
    document.addEventListener('keydown', escape)
} else {
    document.removeEventListener('keydown', escape)
}

function escape(e) {
    if(e.key == 'Escape') {
        kill()
    }
}

function toggle() {
    let theme = localStorage.getItem("theme");
    if(theme && theme == "light") {
        localStorage.removeItem("theme")
        document.documentElement.classList.remove('light')
    } else {
        localStorage.setItem("theme", "light")
        document.documentElement.classList.add('light')
    }
}

let currentView = 'accessibility'
let currentComponent = Accessibility


function navigate(e) {
    currentView = e.detail
    items.forEach(item => {
        item.items.forEach(x => {
            if(x.path == e.detail) {
                currentComponent = x.component
            }
        })
    })
}

let items = [
    {
        title: "user settings",
        items: [
            {
                title: "My Account",
                path: "account",
                component: Account,
            },
            {
                title: "User Profile",
                path: "profile",
                component: Account,
            },
            {
                title: "Privacy & Safety",
                path: "privacy",
                component: Account,
            },
        ]
    },
    {
        title: "app settings",
        items: [
            {
                title: "Appearance",
                path: "appearance",
                component: Appearance,
            },
            {
                title: "Accessibility",
                path: "accessibility",
                component: Accessibility,
            },
            {
                title: "Language",
                path: "language",
                component: Language,
            },
            {
                title: "Notifications",
                path: "notifications",
                component: Account,
            },
        ]
    },
]

</script>

{#if active}
<div class="mask gr-default no-select" 
    on:click|self={kill}>
    <div class="modal gr-center flex flex-column" >
        <div class="con flex">
            <div class="settings-sidebar scrl">
                <Nav 
                items={items}
                on:navigate={navigate}
                currentView={currentView} />
            </div>
            <div class="settings-content scrl">
                <div class="content-container relative">
                    <svelte:component this={currentComponent}/>


                </div>
            </div>
        </div>
    </div>
</div>

<div class="kill fl-co" on:click={kill}>
    <div class="kill-icon gr-center">
        {@html closeBig}
    </div>
    <div class="esc">ESC</div>
</div>

{/if}

<style>
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
    overflow: hidden;
}

.modal {
    background-color: var(--background-2);
    width: 100%;
    height: 100%;
    transition: 0.2s;
    display: grid;
    overflow: hidden;
    z-index: 1001;
}
.con {
    display: grid;
    grid-template-columns: minmax(220px, 35%) 65%;
    width: 100%;
    height: 100%;
    overflow: hidden;
}
.settings-sidebar {
    background-color: var(--background-2);
    display: grid;
    overflow: hidden auto;
}
.settings-content {
    background-color: var(--background-3);
    display: grid;
    overflow: hidden auto;
    padding-left: 2rem;
    padding-right: 1rem;
}
.content-container {
    justify-self: left;
    align-self: center;
    height: 100%;
    padding-top: 4rem;
    max-width: 740px;
    width: 100%;
}
.kill {
    position: fixed;
    right: 1rem;
    top: 1rem;
    z-index: 1003;
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
    scroll-behavior: smooth;
    transition: 0.1s;
}

.scrl:hover {
    scrollbar-color: var(--background-1) transparent;
}

.scrl::-webkit-scrollbar {
    width: 6px;
    border-radius: 1px;
    transition: 0.1s;
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

.kill-icon {
    fill: var(--text);
    border-radius: 50%;
    border: 2px solid var(--text-muted);
    padding: 0.125rem;
    cursor: pointer;
    transition: 0.07s;
}
.kill-icon:hover {
    fill: var(--white);
    border: 2px solid var(--white);
}

.esc {
    color: var(--text-muted);
    font-size: 0.8rem;
    font-weight: bold;
    text-align: center;
    margin-top: 0.5rem;
}
</style>


