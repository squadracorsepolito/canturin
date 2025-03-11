<script lang="ts">
	import {
		BusPanel,
		MessagePanel,
		NetworkPanel,
		NodePanel,
		SignalEnumPanel,
		SignalPanel,
		SignalTypePanel,
		SignalUnitPanel
	} from '.';
	import { getPanelStackState } from './panel-stack-state.svelte';
	import {
		BusIcon,
		MessageIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalIcon,
		SignalTypeIcon,
		SignalUnitIcon
	} from '../components/icon';
	import type { PanelType } from '$lib/state/layout-state.svelte';

	const s = getPanelStackState();

	function getIcon(panelKind: PanelType) {
		switch (panelKind) {
			case 'network':
				return NetworkIcon;
			case 'bus':
				return BusIcon;
			case 'node':
				return NodeIcon;
			case 'message':
				return MessageIcon;
			case 'signal':
				return SignalIcon;
			case 'signal_type':
				return SignalTypeIcon;
			case 'signal_unit':
				return SignalUnitIcon;
			case 'signal_enum':
				return SignalEnumIcon;

			default:
				return NetworkIcon;
		}
	}
</script>

{#snippet viewer()}
	{@const panel = s.stack[s.diplayedIdx]}

	{#if panel.kind === 'network'}
		<NetworkPanel />
	{:else if panel.kind === 'bus'}
		<BusPanel entityId={panel.id} />
	{:else if panel.kind === 'node'}
		<NodePanel entityId={panel.id} />
	{:else if panel.kind === 'message'}
		<MessagePanel entityId={panel.id} />
	{:else if panel.kind === 'signal'}
		<SignalPanel entityId={panel.id} />
	{:else if panel.kind === 'signal_type'}
		<SignalTypePanel entityId={panel.id} />
	{:else if panel.kind === 'signal_unit'}
		<SignalUnitPanel entityId={panel.id} />
	{:else if panel.kind === 'signal_enum'}
		<SignalEnumPanel entityId={panel.id} />
	{/if}
{/snippet}

<div class="flex-1 flex flex-col overflow-hidden">
	<div class="overflow-x-auto">
		<ul class="flex">
			{#each s.stack as panel}
				{@const Icon = getIcon(panel.kind)}

				<li
					class="px-2 py-1 {s.displayedPanel?.id === panel.id
						? 'bg-primary-ghost text-primary'
						: ''}"
				>
					<button class="flex items-center gap-2">
						<span>
							<Icon height={16} width={16} />
						</span>

						<span class="text-sm">
							{panel.name}
						</span>
					</button>
				</li>
			{/each}
		</ul>
	</div>

	<div class="flex-1 overflow-y-auto">
		{#if s.diplayedIdx >= 0}
			{@render viewer()}
		{/if}
	</div>
</div>
