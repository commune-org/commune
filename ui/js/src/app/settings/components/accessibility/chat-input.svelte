<script>
import { store } from '../../../store/store.js'
import Toggle from '../../../components/ui/toggle/toggle.svelte'

let showButton = $store.settings.accessibility.showMessageButton

function toggled() {
    $store.settings.accessibility.showMessageButton = showButton
    store.saveSettings()
}

</script>


<div class="bbt fl-co pb3">  
    <div class="n-t">
        chat input
    </div>


    <div class="mt3 flex">

        <div class="flex-one gr-center">
            Show Send Message button
        </div>

        <div class="flex">
            <Toggle 
            bind:value={showButton} 
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

