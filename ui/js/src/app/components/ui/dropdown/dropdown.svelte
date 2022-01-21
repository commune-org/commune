<script>
import { onMount, onDestroy, createEventDispatcher } from 'svelte'
import { fly } from 'svelte/transition'
import { makeid } from '../../../utils/utils.js'
import { downSmall, checkSmall } from '../../../utils/icons.js'

const dispatch = createEventDispatcher()


export let items;

let selected;

let ready = false;

onMount(() => {
    selected = items?.filter(item => item.default)[0]
    ready = true
})

function select(index) {
    items.forEach((item, i) => {
        if(i == index) {
            item.default = true
        } else {
            item.default = false
        }
    })
    items = items
    selected = items?.filter(item => item.default)[0]
    dispatch('selected', items[index].value)
    kill()
}


let active = false;

function toggle() {
    active = !active
}

$: if(active) {
    setTimeout(() => {
        document.addEventListener('click', killMe)
    }, 10)
} else {
    document.removeEventListener('click', killMe)
}

onDestroy(() => {
    document.removeEventListener('click', killMe)
})

let selector;

let killMe = (e) => {
    if(e.target !== selector && !selector?.contains(e.target)){
        kill()
    }
}

function kill() {
    active = false
}

</script>

{#if ready}
<div class="flex flex-column relative no-select" bind:this={selector}>
    <div class="generic-select flex" on:click={toggle}>
        <div class="gr-center flex-one">
            {selected.text}
        </div>
        <div class="icon gr-center">
            {@html downSmall}
        </div>
    </div>
    {#if active}
        <div class="select-options"
        in:fly="{{ y: -20, duration: 100 }}">
            {#each items as item, i (i)}
                <div class="select-option flex" 
                    class:selected={item.default}
                on:click={select(i)}>
                    <div class="gr-center flex-one">
                        {item.text}
                    </div>
                    {#if item.default}
                        <div class="icon gr-center">
                            {@html checkSmall}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
</div>
{/if}

<style>
.generic-select {
    font-size: 1.2rem;
    background-color: var(--background-2);
    width: 100%;
    padding: 0.5rem 1rem;
    transition: 0.1s;
    border: 1px solid var(--background-1);
    border-radius: 4px;
    cursor: pointer;
}
.select-options {
    position: absolute;
    top: 40px;
    font-size: 1.2rem;
    background-color: var(--background-2);
    width: 100%;
    transition: 0.1s;
    border: 1px solid var(--background-1);
    cursor: pointer;
    z-index: 10000;
}
.select-option {
    padding: 0.5rem 1rem;
}
.select-option:hover {
    background-color: var(--background-3);
}
.selected {
    background-color: var(--background-1);
}
</style>
