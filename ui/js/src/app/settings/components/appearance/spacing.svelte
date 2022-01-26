<script>
import { store } from '../../../store/store.js'
import RangeSlider from '../../../components/ui/range-slider/RangeSlider.svelte'


function update(e) {
    let newVal = e.detail.value
    $store.settings.appearance.eventSpacing = newVal
    store.saveSettings()
}


$: spacing = $store.settings?.appearance?.eventSpacing


</script>

<div class="fl-co">  
    <div class="n-t">
        space between messages & posts
    </div>


    <div class="">

            <RangeSlider 
            min={0}
            max={24}
            range="min"
            step={4}
            values={[spacing]} 
            suffix="px"
            on:change={update}
            def={16}
            pips all="label"/>

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
