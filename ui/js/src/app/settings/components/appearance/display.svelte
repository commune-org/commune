<script>
import { store } from '../../../store/store'
import Select from '../../../components/ui/select/select.svelte'
import {onMount} from 'svelte'


$: md = $store.settings?.appearance?.displayMode

$: isCozy = md && md == "cozy"
$: isCompact = md && md == "compact"

$: modes = [
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


function selected(e) {
    modes = e.detail

    let selected = modes.filter(x => x.selected)[0]
    switch (selected.name) {
        case "compact":
            $store.settings.appearance.displayMode = "compact"
            $store.settings.appearance.eventSpacing = 0
            store.saveSettings()
        break;
        case "cozy":
            $store.settings.appearance.displayMode = "cozy"
            $store.settings.appearance.eventSpacing = 16
            store.saveSettings()
        break;
    }
}


</script>



<div class="fl-co">  
    <div class="n-t">
    message display
    </div>


    <div class="">
        <Select 
            items={modes} 
            on:selected={selected}
        />
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
