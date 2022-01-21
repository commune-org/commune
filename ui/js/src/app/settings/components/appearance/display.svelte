<script>
import { store } from '../../../store/store'
import Select from '../../../components/ui/select/select.svelte'
import {onMount} from 'svelte'

let ready = false;
let modes;

onMount(() => {
    let md = localStorage.getItem("message-display");
    let isCozy = md && md == "cozy"
    let isCompact = md && md == "compact"

    modes = [
        {
            name: "cozy",
            caption: "Cozy: Modern, beautiful, and easy on your eyes.",
            selected: isCozy,
        },
        {
            name: "compact",
            caption: "Compact: Fit more messages on screen at one time. #IRC",
            selected: isCompact,
        },
    ]

    ready = true
})


function selected(e) {
    modes = e.detail

    let selected = modes.filter(x => x.selected)[0]
    switch (selected.name) {
        case "compact":
            localStorage.setItem("message-display", "compact")
            store.updateDisplayMode("compact")
        break;
        case "cozy":
            localStorage.removeItem("message-display")
            store.updateDisplayMode("cozy")
        break;
    }
}

</script>



<div class="fl-co">  
    <div class="n-t">
    message display
    </div>

    <div class="">
        {#if ready}
            <Select 
                items={modes} 
                on:selected={selected}
            />
        {/if}
    </div>
</div>

<style> 
.n-t {
    text-transform: uppercase;
    font-weight: bold;
    font-size: 0.7rem;
    letter-spacing: 1px;
    color: var(--text);
    margin-bottom: 0.5rem;
}
</style>
