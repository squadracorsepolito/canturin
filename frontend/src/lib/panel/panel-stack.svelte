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
	import { getPanelStackState, type Panel, type PanelKind } from './panel-stack-state.svelte';
	import {
		BusIcon,
		CloseIcon,
		MessageIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalIcon,
		SignalTypeIcon,
		SignalUnitIcon
	} from '../components/icon';

	const s = getPanelStackState();

	function getIcon(panelKind: PanelKind) {
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

	function tab(el: HTMLElement, { id }: { id: string }) {
		$effect(() => {
			if (s.displayedPanel?.id === id) {
				el.scrollIntoView({ behavior: 'smooth', block: 'end' });
			}
		});
	}
</script>

{#snippet viewer(panel: Panel)}
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

<div class="flex-1 flex flex-col h-full overflow-hidden">
	<div class="overflow-x-auto border-b-4">
		<ul class="flex">
			{#each s.panels as [_, panel]}
				{@const Icon = getIcon(panel.kind)}
				{@const displayed = s.displayedPanel?.id === panel.id}

				<li
					use:tab={{ id: panel.id }}
					class="group px-2 flex items-center gap-2 {displayed
						? 'bg-primary-ghost text-primary'
						: 'hover:bg-base-content/20'}"
				>
					<button
						onclick={() => s.open(panel.kind, panel.id, panel.name)}
						class="flex items-center gap-2 py-2"
					>
						<span>
							<Icon height={16} width={16} />
						</span>

						<span class="text-sm truncate">
							{panel.name}
						</span>
					</button>

					<button
						onclick={() => s.close(panel.id)}
						class="text-error hover:bg-error-ghost p-1 rounded-btn {!displayed &&
							'invisible'} group-hover:visible"
					>
						<CloseIcon height={16} width={16} />
					</button>
				</li>
			{/each}
		</ul>
	</div>

	<div class="overflow-y-auto flex-1">
		{#if s.displayedPanel}
			{@render viewer(s.displayedPanel)}
		{/if}
	</div>
</div>
