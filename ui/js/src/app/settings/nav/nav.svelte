<script>
import NavItem from './nav-item.svelte'
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

export let items;
export let currentView;

function navigate(e) {
    console.log(e)
}

</script>

<div class="sidebar-nav scrl-sm no-select">
    {#each items as item}
        <div class="nav-section flex-co">
            <div class="n-t">
                {item.title}
            </div>
            {#each item.items as navitem}
                <NavItem 
                on:navigate
                item={navitem} 
                currentView={currentView} />
            {/each}
        </div>
        <div class="nav-sep"></div>
    {/each}
    <div class="logout flex">
        Log Out
    </div>
</div>

<style>
@font-face {
    font-family: Ginto;
    font-weight: 700;
    src: url(/static/fonts/gi700.woff) format("woff");
}
.sidebar-nav {
    justify-self: right;
    align-self: center;
    height: 100%;
    margin-right: 1rem;
    padding-top: 4rem;
    min-width: 200px;
}
.n-t {
    text-transform: uppercase;
    font-family: Ginto;
    font-weight: 700;
    font-size: 0.7rem;
    letter-spacing: 1px;
    color: var(--text-muted);
    margin-left: 0.5rem;
    margin-bottom: 0.5rem;
}
.nav-sep {
    width: 100%;
    height: 1px;
    background-color: var(--background-3);
    margin-top: 0.5rem;
    margin-bottom: 1rem;
}
.logout {
    color: crimson;
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
    margin-bottom: 2px;
}
.logout:hover {
    background-color: var(--background-8);
}
</style>
