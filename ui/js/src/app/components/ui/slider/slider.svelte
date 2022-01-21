<script>
export let range;
export let min = 0;
export let max = 100;
export let step = 1;
export let value;
export let lock;

let w = 50;

$: width = `${w}%`

</script>

<div class="slider-container">
    <div class="bg-track">
    </div>
    <div class="t-track"
        style="--width:{width};">
    </div>
</div>

<style>
:root {
    --bg-track: #4f545c;
    --t-track: red;
}
:root.light {
    --bg-track: #dcddde;
    --t-track: #4f545c;
}

.slider-container {
    position: relative;
    height: 50px;
}
.bg-track {
    border-radius:  500px;
    background-color: var(--bg-track);
    height: 8px;
    width: 100%;
    position: absolute;
    top: 60%;
}
.t-track {
    border-radius:  500px;
    background-color: var(--t-track);
    height: 8px;
    width: var(--width);
    position: absolute;
    top: 60%;
}
</style>
