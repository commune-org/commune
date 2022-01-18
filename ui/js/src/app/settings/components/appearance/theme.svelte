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
            selected: !isLight && !isSync,
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

    let selected = modes.filter(x => x.selected)[0]
    console.log(selected)
    switch (selected.name) {
        case "light":
            localStorage.setItem("theme", "light")
            document.documentElement.classList.add('light')
        break;
        case "dark":
            localStorage.removeItem("theme")
            document.documentElement.classList.remove('light')
        break;
        case "sync":
            localStorage.setItem("theme", "sync")
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.documentElement.classList.remove('light')
            } else { 
                document.documentElement.classList.add('light')
            }
    }
}

</script>



<div class="fl-co">  
    <div class="n-t">
    theme
    </div>

    <div class="">
        {#if ready}
            <Select 
                items={modes} 
                on:selected={selected}
            />
        {/if}
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
