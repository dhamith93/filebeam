<script>
    export let selected;
    export let removeSelectedFile;
</script>

{#if selected.length === 0}
    <div class="box item">
        <p class="filename">No selected files...</p>
    </div>
{/if}

{#each selected as file}
    <div class="box item">
        <p class="filename">{file.Name}</p>
        {#if file.Path !== ''}
            {file.Path}
        {/if}
        <p>
            {`${(file.Size / 1024 / 1024).toFixed(2)} MB`}
        </p>
        <button class="button is-danger" on:click={() => {removeSelectedFile(file)}}>Remove</button>
    </div>
{/each}