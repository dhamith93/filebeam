<script>
    import {notifications} from '../notifications.js'
    import {GetTransfers} from '../../wailsjs/go/main/App.js'
    import {CancelTransfer} from '../../wailsjs/go/main/App.js'
    import {DownloadTransfer} from '../../wailsjs/go/main/App.js'
    import TransferItem from './TransferItem.svelte';
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
        CancelTransfer(ip, filename, isDownload).then(() => {
            notifications.success('Cancelled the file transfer', 3000);
        }).catch(e => {
            notifications.danger(e, 3000);
        });
    };

    const downloadFunc = (host, filename) => {
        DownloadTransfer(host, filename).then(() => {
            notifications.success('Downloading started', 3000);
        }).catch(e => {
            notifications.danger(e, 3000);
        });
    };

    $: {
        setInterval(() => {
            GetTransfers().then((res) => {
                let ts = [];
                res.forEach(t => {
                    let endTime = (t.EndTime === 0) ? Math.floor(Date.now() / 1000) : t.EndTime;
                    let timeDiff = endTime - t.StartTime;
                    let speed = calculateSpeed(timeDiff, t.CompletedBytes);
                    let estimated = (t.SizeBytes - t.CompletedBytes) / speed;
                    ts.push({
                        component: TransferItem,
                        filename: t.File.Name,
                        path: t.File.Path,
                        ip: t.Ip,
                        port: t.FilePort,
                        size: `${(t.SizeBytes / 1024 / 1024).toFixed(2)} MB`,
                        completedSize: `${(t.CompletedBytes / 1024 / 1024).toFixed(2)} MB`,
                        eta: convertToHumanReadableTime(estimated.toFixed(2)),
                        timeSpent: convertToHumanReadableTime(timeDiff.toFixed(2)),
                        speed: (speed > 1000000) ? `${(speed/1024/1024).toFixed(2)} MB/s` : `${(speed/1024).toFixed(2)} KB/s`,
                        done: Math.floor((t.CompletedBytes / t.SizeBytes) * 100),
                        isDownload: t.IsDownload,
                        isCanceled: t.Status === 'cancelled' && t.Status !== 'pending' && t.Status !== 'processing' && t.Status !== 'completes',
                        status: t.Status,
                        cancelFunc: cancelFunc,
                        downloadFunc: downloadFunc
                    })
                });
                transfers = ts;
            }).catch(e => {
                notifications.danger(e, 3000);
            });
        }, 1000);
    }
</script>

<div id="main">
    {#if transfers.length === 0}
        <div class="box item">
            <p class="filename">No transfers...</p>
        </div>
    {/if}
    {#each transfers as transfer}
        <svelte:component this={transfer.component} {...transfer} />
    {/each}
</div>

<style>
    #main {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin: 5px 0 5px 0;
    }
</style>