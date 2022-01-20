<script>
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

export let item;
export let currentView;

$: active = currentView == item.path

function go() {
    dispatch('navigate', item.path)
}
</script>

<div class="nav-item no-select" 
    class:active={active}
    on:click={go}>
    {item.title}
</div>

<style>
.nav-item {
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
    margin-bottom: 2px;
}
.nav-item:hover {
    background-color: var(--background-8);
}
.active {
    background-color: var(--background-11);
}
.active:hover {
    background-color: var(--background-11);
}
</style>
