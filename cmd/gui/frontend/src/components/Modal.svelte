<script>
	export let showTransfers; // boolean

	let dialog; // HTMLDialogElement

	$: if (dialog && showTransfers) dialog.showModal();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog
	bind:this={dialog}
	on:close={() => (showTransfers = false)}
	on:click|self={() => dialog.close()}
>
	<div on:click|stopPropagation>
		<slot />
        <button autofocus on:click={() => dialog.close()}>Close</button>
	</div>
</dialog>

<style>
	dialog {
        background-color: rgba(27, 38, 54, 1);
		width: 600px;
        height: 570px;
		border-radius: 0.2em;
		border: none;
		padding: 0;
        color: #fff;
	}
	dialog::backdrop {
		background: rgba(0, 0, 0, 0.3);
	}
	dialog > div {
		padding: 1em;
	}
	dialog[open] {
		animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
	}
	@keyframes zoom {
		from {
			transform: scale(0.95);
		}
		to {
			transform: scale(1);
		}
	}
	dialog[open]::backdrop {
		animation: fade 0.2s ease-out;
	}
	@keyframes fade {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
    button {
        float: right;
        margin: 5px;
    }
</style>
