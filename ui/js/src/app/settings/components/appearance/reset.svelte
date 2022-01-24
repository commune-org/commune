<script>
import { store } from '../../../store/store.js'
function reset() {
    store.updateFontScaling(16)
    store.updateEventSpacing(16)
    store.updateDisplayMode("cozy")
    store.updateTheme("dark")

    localStorage.removeItem("theme")
    document.documentElement.classList.remove('light')

    localStorage.removeItem("font-scaling")
    localStorage.removeItem("event-spacing")
    localStorage.removeItem("message-display")
}
</script>

<div class="">
    <button on:click={reset}>Reset to Default</button>
</div>

<style>
button {
    background-color: var(--background-2);
    border: 1px solid var(--green);
}
</style>
