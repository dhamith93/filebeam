<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscFolder from "svelte-icons-pack/vsc/VscFolder";
    import VscFileBinary from "svelte-icons-pack/vsc/VscFileBinary";
    import VscArrowLeft from "svelte-icons-pack/vsc/VscArrowLeft";
    import VscArrowRight from "svelte-icons-pack/vsc/VscArrowRight";
    import {GetDirectoryContent} from '../../wailsjs/go/main/App.js'
	export let path;
    export let onFileSelect;
    export let onRefreshPath;
    let selected = [];
    let history = [];
    let content = [];

    export const clearSelection = () => {
        content.forEach(i => {
            if (i.Selected) {
                i.Selected.classList.remove('selected');
            }
        });
        selected = [];
    }

    function refreshPath(p) {
        GetDirectoryContent(p)
        .then(result => {
                content = result;
                selected.forEach(item => {
                    item.Selected.classList.remove('selected');
                })
                selected = [];
                onRefreshPath();
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
    setTimeout(() =>{
        refreshPath(path)
    }, 100);
</script>

<div id="main" class="column">
    <h3>File browser</h3>
    <div id="address-bar">
        <button on:click={handleBackBtnClick}><Icon src="{VscArrowLeft}" /></button> 
        <input type="text" name="path" id="path" bind:value={path} style="--wails-draggable:no-drag"> 
        <button on:click={() => {refreshPath(path)}}><Icon src="{VscArrowRight}" /></button>
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
    }
    #file-tree {
        margin-top: 5px;
        height: calc(100vh - 160px);
        overflow-y: scroll;
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
        margin-bottom: 5px;
        margin-top: 14px;
        padding-left: 5px;
        padding-right: 5px;
    }

    #path {
        font-family: 'JetBrains Mono';
        font-size: 1em;
        width: 100%;
        height: 2rem;
        margin-left: 10px;
        margin-right: 10px;
    }
</style>