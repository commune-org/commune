<script>
import { store } from '../../../store/store.js'
function reset() {
    store.updateFontScaling(16)
    store.updateEventSpacing(16)
    store.updateDisplayMode("cozy")
    store.updateTheme("dark")
    document.documentElement.classList.remove('light')
    store.saveSettings()
}
</script>

<div class="mb4">
    <button on:click={reset}>Reset to Default</button>
</div>

<style>
button {
    background-color: var(--background-2);
    border: 1px solid var(--green);
}
</style>
