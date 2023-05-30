<script>
    import {GetDevices} from '../../wailsjs/go/main/App.js'
    import DeviceListItem from './DeviceListItem.svelte';

    let devices = [];
    function refreshDeviceList() {
        GetDevices().then(result => {
                let newDevices = []
                result.forEach(d => {
                    newDevices.push({
                        host: d,
                        func: () => console.log(d),
                        component: DeviceListItem
                    });
            });
            devices = newDevices;
        })
    }
    refreshDeviceList()
</script>

<div id="devices">
    <button class="btn" on:click={refreshDeviceList}>Refresh</button>
    <div id="device-tree">
        {#each devices as device}
            <svelte:component this={device.component} {...device} />
        {/each}
    </div>
</div>

<style>
    #devices {
        display: flex;
        flex-direction: column;
        width: 30%;
        margin: 20px 60px 0 0;
        align-items: center;
    }
</style>