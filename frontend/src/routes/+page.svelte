<script lang="ts">
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
	import layout from '$lib/state/layout-state.svelte';
</script>

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
	<div>home</div>
{/if}
