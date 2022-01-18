<script>
import {createEventDispatcher} from 'svelte'

const dispatch = createEventDispatcher();

function select() { 
    dispatch("selected", index)
}

export let item;
export let index;
export let multiple;

$: selected = item.selected;

</script>

<div class="s-co flex"
    on:click={() => select()}
    class:selected={item.selected}>
    <div class="gr-center">
        <div class="rad gr-default" 
        class:rad-m={multiple}>
            {#if item.selected}
                <div class="radi gr-center"
                class:radi-m={multiple}>
                </div>
            {/if}
        </div>
    </div>
    <div class="gr-center flex-one ml3 fl-co">
        <div class="">
            {item.caption}
        </div>
    </div>
</div>

<style> 
.s-co {
    border-radius: 5px;
    padding: 0.75rem 1rem;
    background-color: var(--background-2);
    cursor: pointer;
    transition: 0.05s;
    margin-bottom: 0.5rem;
}
.s-co:hover {
    background-color: var(--background-5);
}

.rad {
    border-radius: 50%;
    border: 2px solid var(--white);
    width: 20px;
    height: 20px;
}
.radi {
    border-radius: 50%;
    border: 3px solid var(--background-2);
    background-color: var(--green);
    width: 100%;
    height: 100%;
}

.rad-m {
    border-radius: 5px;
}
.radi-m {
    border-radius: 5px;
}

.selected {
    background-color: var(--background-1);
}
.selected:hover {
    background-color: var(--background-1);
}
</style>
