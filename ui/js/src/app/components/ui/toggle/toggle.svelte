<script>
import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher();

export let value = false;

$: bg = value ? `var(--green)` : `var(--background-4)`

function toggle() {
    value = !value
    dispatch('toggle', value)
}
</script>

<div class="toggle-component"
    style={`--bg: ${bg};`}
    on:click|stopPropagation={toggle}>
    <div class="state"
        class:untoggled={!value}
        class:toggled={value}>
        {#if value}
            <svg viewBox="0 0 20 20" fill="none"><path fill="hsl(139, calc(var(--saturation-factor, 1) * 47.3%), 43.9%)" d="M7.89561 14.8538L6.30462 13.2629L14.3099 5.25755L15.9009 6.84854L7.89561 14.8538Z"></path><path fill="hsl(139, calc(var(--saturation-factor, 1) * 47.3%), 43.9%)" d="M4.08643 11.0903L5.67742 9.49929L9.4485 13.2704L7.85751 14.8614L4.08643 11.0903Z"></path></svg>
        {:else}
            <svg viewBox="0 0 20 20" fill="none"><path fill="hsl(218, calc(var(--saturation-factor, 1) * 4.6%), 46.9%)" d="M5.13231 6.72963L6.7233 5.13864L14.855 13.2704L13.264 14.8614L5.13231 6.72963Z"></path><path fill="hsl(218, calc(var(--saturation-factor, 1) * 4.6%), 46.9%)" d="M13.2704 5.13864L14.8614 6.72963L6.72963 14.8614L5.13864 13.2704L13.2704 5.13864Z"></path></svg>
        {/if}
    </div>
</div>

<style>
.toggle-component {
    width: 40px;
    height: 24px;
    border-radius: 500px;
    background-color: var(--bg);
    cursor: pointer;
    position: relative;
    padding: 3px;
    transition: 0.1s;
}

.state {
    fill: var(--text);
    height: 18px;
    width: 18px;
    position: absolute;
    background-color: var(--white);
    border-radius: 50%;
    transition: 0.1s;
}

.toggled {
    right: 3px;
}
.untoggled {
    right: 19px;
}
</style>
