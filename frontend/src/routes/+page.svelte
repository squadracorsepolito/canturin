<script lang="ts">
	import { IconButton, LinkButton } from '$lib/components/button';
	import { AddIcon, RedoIcon, UndoIcon } from '$lib/components/icon';
	import { Sidebar } from '$lib/components/sidebar';
	import {
		SignalEnumPanel,
		SignalTypePanel,
		BusPanel,
		NodePanel,
		SignalUnitPanel,
		MessagePanel,
		SignalPanel,
		NetworkPanel
	} from '$lib/panel';
	import {
		loadNetwork,
		getNetworkState,
		createNetwork,
		isNetworkLoaded
	} from '$lib/panel/network/state.svelte';
	import { state } from '$lib/state/config-state.svelte';
	import layout from '$lib/state/layout-state.svelte';
	import { Pane, PaneGroup, PaneResizer } from 'paneforge';
	import history from '$lib/state/history-state.svelte';

	function handleUndo() {
		history.undo();
	}

	function handleRedo() {
		history.redo();
	}

	function handleCreateNetwork() {
		createNetwork();
	}

	function handleOpenNetwork(path: string) {
		console.log(path);
		loadNetwork(path);
	}
</script>

{#snippet panel()}
	{#if layout.openPanelType === 'network'}
		<NetworkPanel />
	{:else if layout.openPanelType === 'bus'}
		<BusPanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'node'}
		<NodePanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'message'}
		<MessagePanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'signal'}
		<SignalPanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'signal_type'}
		<SignalTypePanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'signal_unit'}
		<SignalUnitPanel entityId={layout.openPanelId} />
	{:else if layout.openPanelType === 'signal_enum'}
		<SignalEnumPanel entityId={layout.openPanelId} />
	{:else}
		<div>open a panel</div>
	{/if}
{/snippet}

{#snippet editor()}
	<PaneGroup direction="horizontal" class="h-full w-full">
		<Pane defaultSize={15} minSize={5} maxSize={25} class="h-full bg-base-200 flex flex-col">
			<div class="h-12 block bg-base-300 sticky top-0"></div>

			<Sidebar />
		</Pane>

		<PaneResizer
			class="h-full w-1 bg-base-300 data-[active=pointer]:bg-accent hover:bg-accent transition-colors delay-75"
		/>

		<Pane class="flex-1 flex flex-col">
			<div class="h-12 bg-base-200 sticky top-0 block">
				<div class="flex items-center h-full px-5 gap-2">
					<IconButton onclick={handleUndo} disabled={!history.canUndo}>
						<UndoIcon />
					</IconButton>

					<IconButton onclick={handleRedo} disabled={!history.canRedo}>
						<RedoIcon />
					</IconButton>

					{#if !history.history.saved}
						<span>&#42</span>
					{/if}
				</div>
			</div>

			{@render panel()}
		</Pane>
	</PaneGroup>
{/snippet}

{#snippet home()}
	<div class="h-full flex flex-col items-center justify-center gap-20">
		<div>
			<button onclick={handleCreateNetwork} class="btn btn-primary btn-lg">
				<AddIcon />

				<span>Create Network</span>
			</button>
		</div>

		{#if state.cfg.openedNetworks && state.cfg.openedNetworks.length > 0}
			<div>
				<h3 class="pb-5">Opened Networks</h3>

				<ul>
					{#each state.cfg.openedNetworks as openedNet}
						<li>
							<span class="pr-2">{openedNet.name}</span>

							<LinkButton
								label={openedNet.path}
								onclick={() => handleOpenNetwork(openedNet.path)}
							/>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
{/snippet}

{#if isNetworkLoaded()}
	{@render editor()}
{:else}
	{@render home()}
{/if}
