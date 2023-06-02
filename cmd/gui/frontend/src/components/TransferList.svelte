<script>
    import TransferItem from './TransferItem.svelte';

    let transfers = [
        {
            component: TransferItem,
            filename: 'some_file.txt',
            ip: '192.168.111.111',
            size: '7.5Mb',
            completedSize: '6.6Mb',
            eta: '1m 25s',
            speed: '12.4 KB/s',
            done: 40,
            isDownload: false,
            isCanceled: false
        },
        {
            component: TransferItem,
            filename: 'some_file2.txt',
            ip: '192.168.111.111',
            size: '7.5Mb',
            completedSize: '6.6Mb',
            eta: '1m 25s',
            speed: '12.4 KB/s',
            done: 80,
            isDownload: true,
            isCanceled: false
        },
        {
            component: TransferItem,
            filename: 'some_file2.txt',
            ip: '192.168.111.111',
            size: '7.5Mb',
            completedSize: '6.6Mb',
            eta: '1m 25s',
            speed: '12.4 KB/s',
            done: 100,
            isDownload: true,
            isCanceled: false
        },
        {
            component: TransferItem,
            filename: 'some_file2.txt',
            ip: '192.168.111.111',
            size: '7.5Mb',
            completedSize: '6.6Mb',
            eta: '1m 25s',
            speed: '12.4 KB/s',
            done: 80,
            isDownload: true,
            isCanceled: false
        },
        {
            component: TransferItem,
            filename: 'some_file2.txt',
            ip: '192.168.111.111',
            size: '7.5Mb',
            completedSize: '6.6Mb',
            eta: '1m 25s',
            speed: '12.4 KB/s',
            done: 100,
            isDownload: true,
            isCanceled: false
        }
    ];
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