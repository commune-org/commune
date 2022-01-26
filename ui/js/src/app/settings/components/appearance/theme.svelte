<script>
import { store } from '../../../store/store.js'
import Select from '../../../components/ui/select/select.svelte'


$: theme = $store.settings?.appearance?.theme

$: isLight = theme && theme == "light"
$: isBlack = theme && theme == "black"
$: isSync = theme && theme == "sync"

$: modes = [
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
            name: "black",
            caption: "Black",
            selected: isBlack,
        },
        {
            name: "sync",
            caption: "Sync with computer",
            selected: isSync,
        },
    ]

function selected(e) {
    modes = e.detail

    let selected = modes.filter(x => x.selected)[0]
    switch (selected.name) {
        case "light":
            localStorage.setItem("theme", "light")
            document.documentElement.classList.add('light')
            document.documentElement.classList.remove('black')
            store.updateTheme("light")
        break;
        case "black":
            localStorage.setItem("theme", "black")
            document.documentElement.classList.add('black')
            document.documentElement.classList.remove('light')
            store.updateTheme("black")
        break;
        case "dark":
            localStorage.removeItem("theme")
            document.documentElement.classList.remove('light')
            document.documentElement.classList.remove('black')
            store.updateTheme("dark")
        break;
        case "sync":
            localStorage.setItem("theme", "sync")
            store.updateTheme("sync")
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
            <Select 
                items={modes} 
                on:selected={selected}
            />
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
