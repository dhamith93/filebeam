<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import FiFolder from "svelte-icons-pack/fi/FiFolder";
    import FiFilePlus from "svelte-icons-pack/fi/FiFilePlus";
    import {GetDirectoryContent} from '../../wailsjs/go/main/App.js'
	export let path;
    let history = [];
    let content = [];
    let selected = [];

    function refreshPath(p) {
        GetDirectoryContent(p)
        .then(result => {
                content = result;
                selected.forEach(item => {
                    item.Selected.classList.remove('selected');
                })
                selected = [];
            })
            .catch(e => {
                alert('Cannot open given path. Please check the path and try again.');
            });
    }
    function handleFileItemClick(item) {
        if (item.IsDir) {
            history.push(path)
            path = item.Path;
            refreshPath(path);
        } else {
            if (selected.includes(item)) {
                selected.splice(selected.indexOf(item), 1);
                item.Selected.classList.remove('selected');
            } else {
                selected.push(item);
                item.Selected.classList.add('selected');
            }
            console.log(selected)
        }
    }
    function handleBackBtnClick() {
        if (history.length > 0) {
            path = history.pop();
            refreshPath(path);
        }
    }
</script>

<div>
    <button on:click={handleBackBtnClick}>Back</button> <input type="text" name="path" id="path" bind:value={path}> <button on:click={() => {refreshPath(path)}}>Refresh</button>
    <div id="file-tree">
        {#each content as item}
            <div id="file-item" on:click={() => {handleFileItemClick(item)}} on:keydown={() => {handleFileItemClick(item)}} bind:this={item.Selected}>
                {#if item.IsDir}
                    <Icon src="{FiFolder}" />
                {:else}
                    <Icon src="{FiFilePlus}" />
                {/if}
                <p>{item.Name}</p>
    
            </div>
        {/each}
    </div>
</div>


<style>
    div {
        margin: 20px 60px 0 0;
    }
    #file-tree {
        height: 300px;
        overflow: scroll;
    }

    #file-item {
        display: flex;
        margin: 5px;
        padding: 5px;
        cursor: pointer;
        align-items: center;
        overflow: scroll;
    }

    #file-item p {
        margin-top: 0;
        margin-bottom: 0;
        margin-left: 10px;
    }

    #file-item:hover {
        background-color: rgb(86, 86, 86);
    }

    #path {
        min-width: 400px;
    }
</style>