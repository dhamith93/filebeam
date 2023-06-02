<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscArrowDown from "svelte-icons-pack/vsc/VscArrowDown";
    import VscArrowUp from "svelte-icons-pack/vsc/VscArrowUp";
    import VscStopCircle from "svelte-icons-pack/vsc/VscStopCircle";
    import ProgressBar from "@okrad/svelte-progressbar"
    export let transfer = {};
    export let filename;
    export let ip;
    export let completedSize;
    export let size;
    export let eta;
    export let speed;
    export let done;
    export let isDownload;
    export let isCanceled;
    let icnSrc = isDownload ? VscArrowDown : VscArrowUp;
</script>
<div id="item">
    <p class="filename">{filename}</p>
    <p>
        <Icon src="{icnSrc}" /> {ip}
    </p>
    <p>
        {completedSize} / {size} 
    </p>
    <p>
        {eta} | {speed}
    </p>
    {#if done != 100 && !isCanceled}
        <button>Cancel</button>
    {/if}
    {#if isCanceled} 
        <p>Canceled</p>
    {:else}
        <ProgressBar series={[done]} colors={['#EE562E']} addBackground={true} bgColor={"#000"}/>
    {/if}
</div>

<style>
    #item {
        width: 400px;
        background-color: rgb(102, 106, 105);
        margin-bottom: 10px;
    }

    p {
        margin: 2px;
    }
</style>