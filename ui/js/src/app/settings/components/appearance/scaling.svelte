<script>
import DropDown from '../../../components/ui/dropdown/dropdown.svelte'

let items = [
    {
        name: "smallest",
        text: "Smallest",
    },
    {
        name: "smaller",
        text: "Smaller",
    },
    {
        name: "small",
        text: "Small",
    },
    {
        name: "normal",
        text: "Normal",
        default: true,
    },
    {
        name: "large",
        text: "Large",
    },
    {
        name: "larger",
        text: "Larger",
    },
    {
        name: "largest",
        text: "Largest",
    },
]
</script>

<div class="fl-co">  
    <div class="n-t">
        font scaling
    </div>

    <div class="">

        <DropDown items={items} />

    </div>

    ok
    <div class="">
        loll
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
