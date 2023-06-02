<script>
    import {GetHomeDir} from '../../wailsjs/go/main/App.js'
    import FileTree from './FileTree.svelte';
    import DeviceTree from './DeviceTree.svelte';
    let homeDir = '';
    let selected = [];
    GetHomeDir().then(result => homeDir = result);

    function handleFileUpload(host, key) {
        console.log(host, key);
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
    <FileTree path="{homeDir}" onFileSelect={handleFileSelect}/>
    <DeviceTree onFileUpload={handleFileUpload} />
</div>

<style>
    #devices-main {
        display: flex;
        justify-content: center;
    }
</style>