<script>
import Password from './password.svelte'
import Removal from './removal.svelte'
</script>


<div class="account-settings-container fl-co">

    <div class="title">
        <strong>My Account</strong>
    </div>

    <div class="mt3">
        Test
    </div>

    <div class="sp mv4"></div>

    <div class="title">
        Password & Authentication
    </div>

    <div class="mt3">
        <Password />
    </div>

    <div class="sp mv4"></div>

    <div class="">
        <Removal />
    </div>

</div>

<style>
@font-face {
    font-family: Ginto;
    font-weight: 600;
    src: url(/static/fonts/gi600.woff) format("woff");
}
.title {
    font-family: Ginto;
    font-size: 1.2rem;
    color: var(--white);
}
.sp {
    border-bottom: 1px solid var(--background-9);
}
</style>
