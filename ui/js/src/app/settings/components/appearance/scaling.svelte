<script>
import { store } from '../../../store/store.js'
import RangeSlider from '../../../components/ui/range-slider/RangeSlider.svelte'


function update(e) {
    let newVal = e.detail.value
    if(newVal == 16) {
        localStorage.removeItem("font-scaling")
    } else {
        localStorage.setItem("font-scaling", newVal);
    }
    store.updateFontScaling(newVal)
}


$: scaling = $store.settings?.appearance?.fontScaling

</script>

<div class="fl-co">  
    <div class="n-t">
        font scaling
    </div>


    <div class="">

            <RangeSlider 
            min={12}
            max={24}
            range="min"
            step={2}
            values={[scaling]} 
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
