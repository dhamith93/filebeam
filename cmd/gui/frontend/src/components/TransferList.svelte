<script>
    import {GetTransfers} from '../../wailsjs/go/main/App.js'
    import {CancelTransfer} from '../../wailsjs/go/main/App.js'
    import TransferItem from './TransferItem.svelte';

    export let showTransfers;
    let transfers = [];

    const calculateSpeed = (timeInSecs, bytes) => {
        return bytes / timeInSecs;
    };

    const convertToHumanReadableTime = (seconds) => {
        let hours   = Math.floor(seconds / 3600)                 
        let minutes = Math.floor((seconds - (hours * 3600)) / 60)
        seconds = seconds - (hours * 3600) - (minutes * 60)  
        if ( !!hours ) {                                         
            if ( !!minutes ) {                                     
                return `${hours}h ${minutes}m ${seconds.toFixed(2)}s`           
            } else {                                               
                return `${hours}h ${seconds.toFixed(2)}s`                       
            }                                                      
        }                                                        
        if ( !!minutes ) {                                       
            return `${minutes}m ${seconds.toFixed(2)}s`                       
        }                                                        
        return `${seconds.toFixed(2)}s`
    };

    const cancelFunc = (ip, filename, isDownload) => {
        console.log(ip, filename, isDownload)
        CancelTransfer(ip, filename, isDownload).catch(e => {console.log(e)});
    };

    $: if (showTransfers) {
        setInterval(() => {
            GetTransfers().then((res) => {
                let ts = [];
                res.forEach(t => {
                    let endTime = (t.EndTime === '') ? Math.floor(Date.now() / 1000) : t.EndTime;
                    let timeDiff = endTime - t.StartTime;
                    let speed = calculateSpeed(timeDiff, t.CompletedBytes);
                    let estimated = (t.SizeBytes - t.CompletedBytes) / speed;
                    ts.push({
                        component: TransferItem,
                        filename: t.Name,
                        path: t.Path,
                        ip: t.Ip,
                        size: `${(t.SizeBytes / 1024 / 1024).toFixed(2)} MB`,
                        completedSize: `${(t.CompletedBytes / 1024 / 1024).toFixed(2)} MB`,
                        eta: convertToHumanReadableTime(estimated.toFixed(2)),
                        timeSpent: convertToHumanReadableTime(timeDiff.toFixed(2)),
                        speed: (speed > 1000000) ? `${(speed/1024/1024).toFixed(2)} MB/s` : `${(speed/1024).toFixed(2)} KB/s`,
                        done: (t.CompletedBytes / t.SizeBytes) * 100,
                        isDownload: t.IsDownload,
                        isCanceled: t.Status === 'cancelled' || t.Status === 'error',
                        status: t.Status,
                        cancelFunc: cancelFunc
                    })
                });
                transfers = ts;
            }).catch(e => {
                console.log(e);
            });
        }, 1000);
    }
</script>

<div id="main">
    {#each transfers as transfer}
        <svelte:component this={transfer.component} {...transfer} />
    {/each}
</div>

<style>
    #main {
        display: flex;
        flex-direction: column;
        align-items: center;
        height: 500px;
        overflow-y: scroll;
    }
</style>