<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscArrowDown from "svelte-icons-pack/vsc/VscArrowDown";
    import VscArrowUp from "svelte-icons-pack/vsc/VscArrowUp";
    import VscStopCircle from "svelte-icons-pack/vsc/VscStopCircle";
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
<div id="item" class="box">
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
        <button class="button is-danger" on:click={cancelFunc(ip+":"+port, filename, isDownload)}>Cancel</button>
    {/if}
    {#if isDownload && status === 'pending' && !isCanceled}
        <button class="button is-primary" on:click={downloadFunc(ip, filename)}>Download</button>
    {/if}
    {#if !isCanceled}
        <div class="progress-wrapper">
            <progress class="progress is-danger" value="{done}" max="100">{done}%</progress>
        </div>
    {/if}
</div>

<style>
    #item {
        width: 80%;
    }

    .filename {
        font-weight: bold;
    }
    p {
        margin: 2px;
    }

    .progress-wrapper {
        margin: 5px 0 5px 0;
    }
</style>