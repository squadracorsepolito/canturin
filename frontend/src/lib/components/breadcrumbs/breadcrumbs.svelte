<script lang="ts">
	import { EntityKind, type EntityPath } from '$lib/api/canturin';
	import layout from '$lib/state/layout-state.svelte';
	import {
		BusIcon,
		NetworkIcon,
		NodeIcon,
		SignalEnumIcon,
		SignalIcon,
		SignalTypeIcon,
		SignalUnitIcon
	} from '../icon';
	import MessageIcon from '../icon/message-icon.svelte';

	type Props = {
		paths: EntityPath[];
	};

	let { paths }: Props = $props();

	function getIcon(entKind: EntityKind) {
		switch (entKind) {
			case EntityKind.EntityKindNetwork:
				return NetworkIcon;
			case EntityKind.EntityKindBus:
				return BusIcon;
			case EntityKind.EntityKindNode:
				return NodeIcon;
			case EntityKind.EntityKindMessage:
				return MessageIcon;
			case EntityKind.EntityKindSignal:
				return SignalIcon;
			case EntityKind.EntityKindSignalType:
				return SignalTypeIcon;
			case EntityKind.EntityKindSignalUnit:
				return SignalUnitIcon;
			case EntityKind.EntityKindSignalEnum:
				return SignalEnumIcon;
			default:
				return NetworkIcon;
		}
	}
</script>

<div class="breadcrumbs">
	<ul>
		{#each paths as path, idx}
			{@const Icon = getIcon(path.kind)}

			{#if idx < paths.length - 1}
				<li>
					<button
						onclick={() => {
							layout.openPanel0(path.kind, path.entityId);
						}}
						class="underline underline-offset-4 font-medium hover:text-secondary transition-colors
							flex items-center gap-2 text-sm"
					>
						<Icon height={16} width={16}></Icon>

						{path.name}
					</button>
				</li>
			{:else}
				<li class="text-primary">
					<Icon height={16} width={16} />

					<span class="text-sm font-medium pl-2">
						{path.name}
					</span>
				</li>
			{/if}
		{/each}
	</ul>
</div>
