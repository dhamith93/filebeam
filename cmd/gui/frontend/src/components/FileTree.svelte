<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscFolder from "svelte-icons-pack/vsc/VscFolder";
    import VscFileBinary from "svelte-icons-pack/vsc/VscFileBinary";
    import VscArrowLeft from "svelte-icons-pack/vsc/VscArrowLeft";
    import VscArrowRight from "svelte-icons-pack/vsc/VscArrowRight";
    import {GetDirectoryContent} from '../../wailsjs/go/main/App.js'
	export let path;
    export let onFileSelect;
    let selected = [];
    let history = [];
    let content = [];

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
            onFileSelect(item)
            if (selected.includes(item)) {
                selected.splice(selected.indexOf(item), 1);
                item.Selected.classList.remove('selected');
            } else {
                selected.push(item);
                item.Selected.classList.add('selected');
            }
        }
    }
    function handleBackBtnClick() {
        if (history.length > 0) {
            path = history.pop();
            refreshPath(path);
        }
    }
</script>

<div id="main">
    <h3>File browser</h3>
    <div id="address-bar">
        <button on:click={handleBackBtnClick}><Icon src="{VscArrowLeft}" /></button> <input type="text" name="path" id="path" bind:value={path}> <button on:click={() => {refreshPath(path)}}><Icon src="{VscArrowRight}" /></button>
    </div>
    <div id="file-tree">
        {#each content as item}
            <div id="file-item" on:click={() => {handleFileItemClick(item)}} on:keydown={() => {handleFileItemClick(item)}} bind:this={item.Selected}>
                <span>
                    {#if item.IsDir}
                        <Icon src="{VscFolder}" />
                    {:else}
                        <Icon src="{VscFileBinary}" />
                    {/if}
                </span>
                <p>{item.Name}</p>
    
            </div>
        {/each}
    </div>
</div>


<style>
    #main {
        padding: 5px;
        border: 3px solid rgb(122, 122, 122);
        height: 500px;
    }
    #file-tree {
        max-width: 500px;
        height: 400px;
        overflow: scroll;
    }

    #file-item {
        display: flex;
        margin: 5px;
        padding: 5px;
        cursor: pointer;
        align-items: center;
    }

    #file-item p {
        margin-top: 0;
        margin-bottom: 0;
        margin-left: 10px;
        overflow: hidden;
        text-overflow: ellipsis;
        text-align: left;
    }

    #file-item:hover {
        background-color: rgb(86, 86, 86);
    }

    #address-bar {
        display: flex;
        align-items: center;
    }

    #path {
        min-width: 400px;
        height: 1rem;
        margin-left: 10px;
        margin-right: 10px;
    }
</style>