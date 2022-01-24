<script>
import RangeSlider from '../../../components/ui/range-slider/RangeSlider.svelte'

</script>

<div class="fl-co">  
    <div class="n-t">
        saturation
    </div>


    <div class="co">
        Reduce the saturation of colors within the application, for those with color sensitivities. This does not affect the saturation of images, videos, role colors or other user-provided content by default.
    </div>


    <div class="mt2">
            <RangeSlider 
            min={0}
            max={100}
            range="min"
            step={10}
            values={[100]} 
            suffix="%"
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
.co {
    color: var(--text-alt);
    line-height: 1.3rem;
}
</style>

