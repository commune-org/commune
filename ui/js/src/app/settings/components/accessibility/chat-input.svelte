<script>
import Toggle from '../../../components/ui/toggle/toggle.svelte'

let showButton = false;

function toggled() {
    console.log(showButton)
}

</script>


<div class="fl-co">  
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
</style>

