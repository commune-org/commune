<script>
import Saturation from './saturation.svelte'
import ChatInput from './chat-input.svelte'
import LimitWidth from './limit-width.svelte'
</script>

<div class="accessibility-settings-container fl-co">

    <div class="title">
        Accessibility
    </div>
    
    <div class="mt4">
        <Saturation/>
    </div>

    <div class="mt4">
        <LimitWidth/>
    </div>

    <div class="mt4">
        <ChatInput/>
    </div>



</div>

<style>
.title {
    font-size: 1.2rem;
    font-weight: bold;
    color: var(--white);
}
</style>

