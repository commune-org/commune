<script>
import { store } from '../../../store/store.js'
import{ onMount } from 'svelte'
import tippy from 'tippy.js'
import 'tippy.js/animations/perspective.css';

export let server;
export let room;

$: server = store.getServer(`/${server}`)

onMount(() => {
    initMenu()
})

let menuActive = false;

let container;
let menu;
let content;
let initMenu = () => {
    menu = tippy(container, {
        content: content.innerHTML,
        theme: 'menu',
        placement: 'bottom',
        trigger: 'click',
        animation: 'perspective',
        interactive: true,
        arrow: false,
        delay: 0,
        duration: 70,
        allowHTML: true,
        onShown(menu) {
            menuActive = true
            console.log("heheh", menuActive)
        },
        onHide(menu) {
            menuActive = false
        },
    });
}

</script>

<template bind:this={content}>
    <div class="server-menu fl-co">
        <div class="server-menu-item">
            One
        </div>
    </div>
</template>

<div class="sp-he gr-default h-100" 
    class:menu-active={menuActive}
    bind:this={container}>
    <div class="gr-center w-100 ph3 flex">
        <div class="name clmp-1 flex-one">
            {server?.name}
        </div>
        <div class="">
            {#if menuActive}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"></path></svg>
            {:else}
                <svg viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M12.78 6.22a.75.75 0 010 1.06l-4.25 4.25a.75.75 0 01-1.06 0L3.22 7.28a.75.75 0 011.06-1.06L8 9.94l3.72-3.72a.75.75 0 011.06 0z"></path></svg>
            {/if}
        </div>
    </div>
</div>

<style>
.sp-he {
    transition: 0.2s;
    cursor: pointer;
    fill: var(--text);
}
.sp-he:hover {
    background-color: var(--background-6);
}
.menu-active {
    background-color: var(--background-6);
}
.server-menu{
    width: 206px;
}
.server-menu-item{
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
}
.server-menu-item:hover{
    color: white;
    background-color: var(--secondary);
}
.name {
    font-family: Ginto;
    font-weight: 600;
}

</style>
