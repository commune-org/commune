<script>
import RangeSlider from '../../../components/ui/range-slider/RangeSlider.svelte'


function deal(e) {
    console.log(e.detail)
}

</script>

<div class="fl-co">  
    <div class="n-t">
        font scaling
    </div>

    <div class="">

        <RangeSlider 
        min={12}
        max={24}
        step={2}
        values={[16]} 
        suffix="px"
        on:change={deal}
        pips all='label'/>


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
