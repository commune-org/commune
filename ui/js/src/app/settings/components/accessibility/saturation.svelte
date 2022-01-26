<script>
import { store } from '../../../store/store.js'
import RangeSlider from '../../../components/ui/range-slider/RangeSlider.svelte'

function update(e) {
    let newVal = e.detail.value
    let v = newVal/100
    if(newVal == 100) {
        document.documentElement.style.removeProperty('filter')
    } else {
        document.documentElement.style.filter = `saturate(${v})`
    }
    $store.settings.accessibility.saturation = v
    store.saveSettings()
}

$: saturation = $store.settings?.accessibility?.saturation * 100

</script>

<div class="fl-co">  
    <div class="n-t">
        saturation
        {saturation}
    </div>


    <div class="co">
        Reduce the saturation of colors throughout the entire application.
    </div>


    <div class="mt2">
            <RangeSlider 
            min={0}
            max={100}
            range="min"
            step={10}
            values={[saturation]} 
            suffix="%"
            on:change={update}
            def={100}
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
.co {
    color: var(--text-alt);
    line-height: 1.3rem;
}
</style>

