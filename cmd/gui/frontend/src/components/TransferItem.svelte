<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscArrowDown from "svelte-icons-pack/vsc/VscArrowDown";
    import VscArrowUp from "svelte-icons-pack/vsc/VscArrowUp";
    import VscStopCircle from "svelte-icons-pack/vsc/VscStopCircle";
    import ProgressBar from "@okrad/svelte-progressbar"
    export let filename;
    export let path;
    export let ip;
    export let port;
    export let completedSize;
    export let size;
    export let eta;
    export let speed;
    export let done;
    export let isDownload;
    export let isCanceled;
    export let status;
    export let timeSpent;
    export let cancelFunc;
    export let downloadFunc;
    let icnSrc = isDownload ? VscArrowDown : VscArrowUp;
</script>
<div id="item">
    <p class="filename">{filename}</p>
    {#if path !== ''}
        {path}
    {/if}
    <p>
        <Icon src="{icnSrc}" /> {ip}
    </p>
    <p>
        {completedSize} / {size} 
    </p>
    <p>
        {#if !isCanceled && status !== 'pending'}
            ETA: {eta} | Spent: {timeSpent}
        {/if}
    </p>
    <p>
        {speed}
    </p>
    <p>
        {status.toUpperCase()}
    </p>
    {#if done != 100 && !isCanceled}
        <button on:click={cancelFunc(ip+":"+port, filename, isDownload)}>Cancel</button>
    {/if}
    {#if isDownload && status === 'pending' && !isCanceled}
        <button on:click={downloadFunc(ip, filename)}>Download</button>
    {/if}
    {#if !isCanceled} 
        <ProgressBar series={[done]} colors={['#EE562E']} addBackground={true} bgColor={"#000"}/>
    {/if}
</div>

<style>
    #item {
        width: 450px;
        background-color: rgb(102, 106, 105);
        margin-bottom: 10px;
    }

    .filename {
        font-weight: bold;
    }
    p {
        margin: 2px;
    }
</style>