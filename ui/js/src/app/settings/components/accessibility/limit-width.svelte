<script>
import { store } from '../../../store/store.js'
import Toggle from '../../../components/ui/toggle/toggle.svelte'

let limitDisplayWidth = $store.settings.accessibility.limitDisplayWidth

function toggled() {
    $store.settings.accessibility.limitDisplayWidth = limitDisplayWidth
    store.saveSettings()
}

</script>


<div class="bbt fl-co pb3">  
    <div class="n-t">
        limit display width
    </div>


    <div class="mt3 flex">

        <div class="flex-one gr-center">
            Limit application width on large displays.
        </div>

        <div class="flex">
            <Toggle 
            bind:value={limitDisplayWidth} 
            on:toggle={toggled}
            />
        </div>

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
.bbt {
    border-bottom: 1px solid var(--background-9);
}
</style>

