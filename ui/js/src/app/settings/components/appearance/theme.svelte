<script>
import Select from '../../../components/ui/select/select.svelte'
import {onMount} from 'svelte'

let ready = false;
let modes;

onMount(() => {
    let theme = localStorage.getItem("theme");
    let isLight = theme && theme == "light"
    let isSync = theme && theme == "sync"

    modes = [
        {
            name: "dark",
            caption: "Dark",
            selected: !isLight,
        },
        {
            name: "light",
            caption: "Light",
            selected: isLight,
        },
        {
            name: "sync",
            caption: "Sync with computer",
            selected: isSync,
        },
    ]



    ready = true
})


function selected(e) {
    modes = e.detail

    console.log(modes)
}

</script>

<div class="">  
    {#if ready}
        <Select 
            items={modes} 
            on:selected={selected}
        />
    {/if}
</div>
