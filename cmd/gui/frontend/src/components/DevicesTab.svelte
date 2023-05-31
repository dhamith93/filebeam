<script>
    import {GetHomeDir} from '../../wailsjs/go/main/App.js'
    import FileTree from './FileTree.svelte';
    import DeviceTree from './DeviceTree.svelte';
    let homeDir = '';
    let deviceKey = '';
    GetHomeDir().then(result => homeDir = result);
</script>
<div id="devices-main">
    <FileTree path="{homeDir}" />
    <DeviceTree />
</div>
<div id="transfers">

</div>

<style>
    #devices-main {
        display: flex;
        justify-content: center;
    }
</style>