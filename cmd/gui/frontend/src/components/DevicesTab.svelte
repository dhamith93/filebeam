<script>
    import {GetHomeDir} from '../../wailsjs/go/main/App.js'
    import {AddToQueue} from '../../wailsjs/go/main/App.js'
    import FileTree from './FileTree.svelte';
    import DeviceTree from './DeviceTree.svelte';
    let homeDir = '';
    let selected = [];
    let fileTree;
    GetHomeDir().then(result => homeDir = result);

    function handleFileUpload(host, key) {
        AddToQueue(selected, host, key).then(() => {
            fileTree.clearSelection();
            selected = [];
        }).catch(e => {
            console.log(e);
        });
    }

    function handleFileSelect(item) {
        if (selected.includes(item)) {
            selected.splice(selected.indexOf(item), 1);
        } else {
            selected.push(item);
        }
        console.log(selected)
    }
</script>
<div id="devices-main">
    <FileTree path="{homeDir}" onFileSelect={handleFileSelect} bind:this={fileTree}/>
    <DeviceTree onFileUpload={handleFileUpload} />
</div>

<style>
    #devices-main {
        display: flex;
        justify-content: center;
    }
</style>